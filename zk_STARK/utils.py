from merkle_tree import hash
import os
import json

# Get the set of powers of R, until but not including when the powers
# loop back to 1


def get_power_cycle(r, modulus):
    o = [1, r]
    while o[-1] != 1:
        o.append((o[-1] * r) % modulus)
    return o[:-1]

# Extract pseudorandom indices from entropy


def get_pseudorandom_indices(seed, modulus, count, exclude_multiples_of=0):
    assert modulus < 2**24
    data = seed
    while len(data) < 4 * count:
        data += hash(data[-32:])
    if exclude_multiples_of == 0:
        return [int.from_bytes(data[i: i+4], 'big') % modulus for i in range(0, count * 4, 4)]
    else:
        real_modulus = modulus * \
            (exclude_multiples_of - 1) // exclude_multiples_of
        o = [int.from_bytes(data[i: i+4], 'big') %
             real_modulus for i in range(0, count * 4, 4)]
        return [x+1+x//(exclude_multiples_of-1) for x in o]


def is_a_power_of_2(x):
    return True if x == 1 else False if x % 2 else is_a_power_of_2(x//2)


def extend_user_data(ids, coins, uts):
    user_num = len(ids)
    coins_num = len(coins)
    for i in range(coins_num):
        assert (len(coins[i]) == user_num)

    extended_ids = [0] * uts
    extended_coins = []
    for i in range(user_num):
        extended_ids = extended_ids + ([0]*(uts-2) + [ids[i]] + [0])

    for i in range(coins_num):
        balances = [0] * uts
        for j in range(user_num):
            balances = balances + ([0]*(uts-2) + [coins[i][j]] + [0])
        extended_coins.append(balances)
    del balances
    return extended_ids, extended_coins


def get_sum_trace(coins, uts, modulus):
    user_num = len(coins[0])//uts

    trace = [0]*len(coins[0])
    sum_values = [0]*(len(coins)+1)

    # calculate sum_values for each coins and trace[uts * i + uts - 2] for each user
    for i in range(user_num):
        for j in range(len(coins)):
            sum_values[j] = (sum_values[j] + coins[j]
                             [i * uts + uts - 2]) % modulus
            trace[uts * i + uts - 2] = (trace[uts * i + uts - 2] +
                                        coins[j][i * uts + uts - 2]) % modulus

    # calcualte total sum_values
    for i in range(len(coins)):
        sum_values[-1] += sum_values[i]

    # calculate trace[uts * i + uts - 1]
    for i in range(1, user_num):
        trace[uts * i + uts - 1] = (trace[uts * (i - 1) +
                                    uts - 1] + trace[uts * i + uts - 2]) % modulus

    # calculate trace[i] when i % uts != {uts-2, uts-1}
    for i in range(user_num):
        for j in range(uts-3):
            trace[uts*i+uts-3-j] = trace[uts*i+uts-2-j] // 4

    # calculate coins sum_values trace
    for i in range(len(coins)):
        for j in range(1, user_num):
            coins[i][uts * j + uts - 1] = (
                coins[i][uts * (j - 1) + uts - 1] + coins[i][uts * j + uts - 2]) % modulus

    return trace, sum_values, coins


def pad(ids, coins, max_user_num):
    padded_num = max_user_num
    user_num = len(ids)
    while padded_num > 2 * (user_num + 1):
        padded_num //= 2
    ids = ids + [0] * (padded_num - user_num - 1)
    for i in range(len(coins)):
        coins[i] = coins[i] + [0] * (padded_num - user_num - 1)
    return ids, coins


def get_entry_data(entry_data, index):
    x = b''
    for j in range(len(entry_data)):
        if (type(entry_data[j][0]) == list):
            for k in range(len(entry_data[j])):
                x = x + entry_data[j][k][index].to_bytes(32, 'big')
        else:
            x = x + entry_data[j][index].to_bytes(32, 'big')
    return x


def get_entries(array):
    entries_len = len(array[0][0]) if type(
        array[0][0]) == list else len(array[0])
    entries = []
    for i in range(entries_len):
        x = get_entry_data(array, i)
        entries.append(x)
    del x, entries_len
    return entries


def get_leaves(array):
    leaves_len = len(array[0][0]) if type(
        array[0][0]) == list else len(array[0])
    leaves = []
    for i in range(leaves_len):
        x = b''
        for j in range(len(array)):
            if (type(array[j][0]) == list):
                temp = b''
                for k in range(len(array[j])):
                    temp = temp + array[j][k][i].to_bytes(32, 'big')
                x = x + hash(temp)
            else:
                x = x + array[j][i].to_bytes(32, 'big')
        leaves.append(hash(x))
    del x, leaves_len
    return leaves


def calculate_l(seed, powers, array, modulus):
    l_len = len(array[0][0]) if type(array[0][0]) == list else len(array[0])
    k_len = 0
    l = []
    for i in range(len(array)):
        if type(array[i][0]) == int:
            k_len += 2
        else:
            k_len += 2 * len(array[i])
    k = [int.from_bytes(hash(seed + i.to_bytes(32, 'big')),
                        'big') % modulus for i in range(k_len)]
    index = 0
    for i in range(l_len):
        x = 0
        for j in range(len(array)):
            if (type(array[j][0]) == list):
                for m in range(len(array[j])):
                    x = (x + k[index] * array[j][m][i] % modulus) % modulus
                    x = (x + (k[index+1] * array[j][m][i] %
                         modulus) * powers[i] % modulus) % modulus
                    index = (index + 2) % k_len
            else:
                x = (x + (k[index] * array[j][i] % modulus)) % modulus
                x = (x + (k[index+1] * array[j][i] % modulus)
                     * powers[i] % modulus) % modulus
                index = (index + 2) % k_len
        l.append(x)
    del x, index, l_len, k_len
    return l


def verify_l(k, power, l, array, extra, modulus):
    index = 0
    for i in range(len(array)):
        if (type(array[i]) == list):
            for j in range(len(array[i])):
                l = (l - k[index] * array[i][j] % modulus) % modulus
                l = (l - (k[index+1] * array[i][j] %
                     modulus) * power % modulus) % modulus
                index = index + 2

        else:
            l = (l - k[index] * array[i] % modulus) % modulus
            l = (l - (k[index+1] * array[i] %
                 modulus) * power % modulus) % modulus
            index = index + 2

    return (l - extra) % modulus == 0


def check_entry_hash(main_branch_leaves, data, coins_num, modulus):
    data_length = len(data[0])
    for i in range(len(main_branch_leaves)//4):
        user_random = hash(data[4*i][data_length-2*coins_num*32-5*32:])
        b_hash = hash(data[4*i][32:data_length-2*coins_num*32-6*32])
        assert main_branch_leaves[4*i] == hash(data[4*i][:32] + b_hash + data[4*i]
                                               [data_length-2*coins_num*32-6*32:data_length-2*coins_num*32-5*32] + user_random)
        assert main_branch_leaves[4*i+1] == hash(
            data[4*i+1][0] + data[4*i+1][1] + data[4*i+1][2] + data[4*i+1][3])
        assert main_branch_leaves[4*i+2] == hash(
            data[4*i+2][0] + hash(data[4*i+2][1]) + data[4*i+2][2] + data[4*i+2][3])
        assert main_branch_leaves[4*i+3] == hash(
            data[4*i+3][0] + hash(data[4*i+3][1]) + data[4*i+3][2] + data[4*i+3][3])
    return


def save_mtree_entries_data(data_path, mtree_entries_data):
    if not os.path.exists(data_path):
        os.mkdir(data_path)

    with open(data_path + "mtree_entries_data.json", "w") as ff:
        mtree_entries_data_json = {}
        for i in range(len(mtree_entries_data)):
            mtree_entries_data_json[str(i)] = str(mtree_entries_data[i].hex())
        json.dump(mtree_entries_data_json, ff)
    return


def save_data(data_path, sum_proof, mtree, sum_values, coins):
    if not os.path.exists(data_path):
        os.mkdir(data_path)

    with open(data_path + "sum_proof.json", "w") as ff:
        sum_proof_json = {
            "steps": sum_proof[0],
            "uts": sum_proof[1],
            "mtree_root": sum_proof[2].hex(),
            "l_mtree_root": sum_proof[3].hex(),
            "pow_nonce": sum_proof[4].hex(),
            "mtree_branches": bytes_array_to_hex(sum_proof[5]),
            "mtree_entries_data": bytes_array_to_hex(sum_proof[6]),
            "l_mtree_branches": bytes_array_to_hex(sum_proof[7]),
            "low_degree_proof": bytes_array_to_hex(sum_proof[8])
        }
        json.dump(sum_proof_json, ff)

    with open(data_path + "mtree.json", "w") as ff:
        mtree = bytes_array_to_hex(mtree)
        mtree_json = {
            "mtree": mtree
        }
        json.dump(mtree_json, ff)

    with open(data_path + "sum_values.json", "w") as ff:
        coins_json = {}
        j = 0
        for coin in coins:
            coins_json[coin] = sum_values[j]
            j += 1
        coins_json["total_value"] = sum_values[-1]
        json.dump(coins_json, ff)

    return


def transform_into_field(x, modulus):
    for i in range(len(x)):
        if type(x[i]) == list:
            for j in range(len(x[i])):
                x[i][j] = x[i][j] % modulus
        else:
            x[i] = x[i] % modulus
    return


def check_sum_values(sum_values, modulus):
    cal_sum = 0
    for i in range(len(sum_values)-1):
        assert sum_values[i] >= 0 and sum_values[i] < modulus, "invalid sum amounts"
        cal_sum += sum_values[i]
    assert cal_sum == sum_values[-1]
    return


def bytes_array_to_hex(array):
    for i in range(len(array)):
        if type(array[i]) == list:
            array[i] = bytes_array_to_hex(array[i])
        elif type(array[i]) == bytes:
            array[i] = array[i].hex()
    return array


def hex_array_to_bytes(array):
    for i in range(len(array)):
        if type(array[i]) == list:
            array[i] = hex_array_to_bytes(array[i])
        elif type(array[i]) == str:
            array[i] = bytes.fromhex(array[i])
    return array


def proof_of_work(seed, pow_bits):
    nonce = int.from_bytes(seed, 'big')
    while True:
        x = int.from_bytes(hash(seed + nonce.to_bytes(32, 'big')), 'big')
        if x >> (256 - pow_bits) == 0:
            break
        nonce += 1
    return nonce.to_bytes(32, 'big')
