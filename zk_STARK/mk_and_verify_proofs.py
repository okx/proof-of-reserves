from por_stark import mk_por_proof, verify_por_proof
from permuted_tree import mk_branch, verify_branch, keccak_256
from poly_utils import PrimeField
from constants import *
from utils import check_entry_hash, hex_array_to_bytes

import random
import time
import json
import os
import gc
import re

f = PrimeField(MODULUS)


def init_user_data(batch_size, batch_index, config_path):
    with open(config_path, "r") as ff:
        config_json = json.load(ff)
        coins = config_json["coins"]
    data = []
    coins_len = len(coins)
    for i in range(batch_size):
        items = {"id": str(keccak_256(i.to_bytes(32, 'big')).hex())}
        for coin in coins:
            items[coin] = str(random.randrange(4**(UTS16 - 2))//coins_len)
        data.append(items)
    with open(USER_DATA_PATH + "batch" + str(batch_index) + ".json", "w") as f:
        json.dump(data, f)
    return

# path: json path of user data


def read_user_data(data_path, config_path):
    with open(data_path, "r") as f:
        data = json.load(f)
    with open(config_path, "r") as ff:
        config_json = json.load(ff)
        coins = config_json["coins"]
    ids = []
    balances = [[]]*(len(coins))
    for item in data:
        ids.append(int(item["id"], 16))
        j = 0
        for coin_name in coins:
            if coin_name in item.keys():
                balances[j] = balances[j] + [int(item[coin_name], 10)]
            else:
                balances[j] = balances[j] + [0]
            j += 1
    return ids, balances

# input_path: json path of user data, xxx/xxx.json
# output_path: batch path of proof data, xxx/xxx/a1/


def mk_batch_proof(uts, input_path, output_path, config_path):
    ids, balances = read_user_data(input_path, config_path)
    assert len(ids) < MAX_USER_NUM_FOR_ONE_BATCH, "too much users in one batch"
    with open(config_path, "r") as ff:
        config_json = json.load(ff)
        main_coins_num = int(config_json["main_coins_num"])
        coins = config_json["coins"]
    mk_por_proof(ids, balances, uts, output_path, main_coins_num, coins)
    return

# batch path of proof data, xxx/xxx/


def verify_batch_proof(input_path, config):
    main_coins_num = int(config["main_coins_num"])
    coins = config["coins"]
    with open(input_path + "sum_proof.json", "r") as ff:
        sum_proof_json = json.load(ff)
        sum_proof = [sum_proof_json["steps"],
                     sum_proof_json["uts"],
                     bytes.fromhex(sum_proof_json["mtree_root"]),
                     bytes.fromhex(sum_proof_json["l_mtree_root"]),
                     bytes.fromhex(sum_proof_json["pow_nonce"]),
                     hex_array_to_bytes(sum_proof_json["mtree_branches"]),
                     hex_array_to_bytes(sum_proof_json["mtree_entries_data"]),
                     hex_array_to_bytes(sum_proof_json["l_mtree_branches"]),
                     hex_array_to_bytes(sum_proof_json["low_degree_proof"])]
    with open(input_path + "sum_values.json", "r") as ff:
        sum_values_json = json.load(ff)
        sum_values = []
        for coin in coins:
            sum_values.append(sum_values_json[coin])
        sum_values.append(sum_values_json["total_value"])

    assert verify_por_proof(sum_values, sum_proof,
                            main_coins_num), "invalid batch proof"
    return sum_values

# input_path: basic batch path, xxx/batches/
# output_path: trunk path, xxx/xxx/


def mk_trunk_proof(input_path, output_path, config_path):
    with open(config_path, "r") as ff:
        config_json = json.load(ff)
        main_coins_num = int(config_json["main_coins_num"])
        coins = config_json["coins"]
    ids = []
    values = [[]]*(len(coins))

    a_count = 0
    b_count = 0
    for root, dirs, files in os.walk(input_path):
        for dir in dirs:
            if dir[0] == 'a':
                a_count += 1
            elif dir[0] == 'b':
                b_count += 1

    for i in range(a_count):
        with open(input_path + "a" + str(i) + "/mtree.json", "r") as ff:
            ids.append(int(json.load(ff)["mtree"][1], 16))
        with open(input_path + "a" + str(i) + "/sum_values.json", "r") as ff:
            sum_values = json.load(ff)
            j = 0
            for coin in coins:
                values[j] = values[j] + [sum_values[coin]]
                j += 1

    for i in range(b_count):
        with open(input_path + "b" + str(i) + "/mtree.json", "r") as ff:
            ids.append(int(json.load(ff)["mtree"][1], 16))
        with open(input_path + "b" + str(i) + "/sum_values.json", "r") as ff:
            sum_values = json.load(ff)
            j = 0
            for coin in coins:
                values[j] = values[j] + [sum_values[coin]]
                j += 1

    mk_por_proof(ids, values, UTS_FOR_TRUNK,
                 output_path, main_coins_num, coins)
    return

# input_path: trunk path, xxx/xxx/


def verify_trunk_proof(input_path, config):
    coins = config["coins"]
    main_coins_num = int(config["main_coins_num"])
    with open(input_path + "sum_proof.json", "r") as ff:
        sum_proof_json = json.load(ff)
        sum_proof = [sum_proof_json["steps"],
                     sum_proof_json["uts"],
                     bytes.fromhex(sum_proof_json["mtree_root"]),
                     bytes.fromhex(sum_proof_json["l_mtree_root"]),
                     bytes.fromhex(sum_proof_json["pow_nonce"]),
                     hex_array_to_bytes(sum_proof_json["mtree_branches"]),
                     hex_array_to_bytes(sum_proof_json["mtree_entries_data"]),
                     hex_array_to_bytes(sum_proof_json["l_mtree_branches"]),
                     hex_array_to_bytes(sum_proof_json["low_degree_proof"])]

    with open(input_path + "sum_values.json", "r") as ff:
        sum_values_json = json.load(ff)
        sum_values = []
        for coin in coins:
            sum_values.append(sum_values_json[coin])
        sum_values.append(sum_values_json["total_value"])
    assert verify_por_proof(sum_values, sum_proof,
                            main_coins_num), "invalid trunk proof"
    return sum_values

# batch_index: batch index in trunk
# input_batch_path: batch path, xxx/batches/a1/
# input_trunk_path: trunk path, xxx/trunk/
# output_path: path for saving inclusion data, xxx/inclusion_proof_data/a1/


def mk_inclusion_proof(batch_index, uts, input_batch_path, input_trunk_path, output_path, config_path):
    start_time = time.time()
    with open(config_path, "r") as ff:
        config_json = json.load(ff)
        coins = config_json["coins"]
    coin_num = len(coins)

    with open(input_trunk_path + "mtree.json", "r") as ff:
        trunk_mtree = json.load(ff)["mtree"]
    with open(input_trunk_path + "mtree_entries_data.json", "r") as ff:
        trunk_mtree_entries_data = json.load(ff)

    batch_entry_data = bytes.fromhex(
        trunk_mtree_entries_data[str(batch_index)])
    del trunk_mtree_entries_data
    gc.collect()

    trunk_inclusion_proof = {}
    trunk_inclusion_proof["trunk_mtree_root"] = trunk_mtree[1]
    trunk_inclusion_proof["batch_id"] = str(
        batch_entry_data[(coin_num+1)*32:(coin_num+2)*32].hex())
    trunk_inclusion_proof["total_value"] = str(
        int.from_bytes(batch_entry_data[:32], 'big'))
    j = 0
    for coin in coins:
        value = int.from_bytes(batch_entry_data[(j+1)*32:(j+2)*32], 'big')
        if value > MAX_USER_VALUE:
            value = value - MODULUS
        trunk_inclusion_proof[coin] = str(value)
        j += 1
    trunk_inclusion_proof["random_number"] = str(
        batch_entry_data[len(batch_entry_data)-32:].hex())
    trunk_inclusion_proof["merkle_path"] = mk_branch(
        trunk_mtree, (UTS_FOR_TRUNK * (batch_index + 1) + UTS_FOR_TRUNK-2) * EXTENSION_FACTOR)

    with open(input_batch_path + "mtree.json", "r") as ff:
        batch_mtree = json.load(ff)["mtree"]
    with open(input_batch_path + "mtree_entries_data.json", "r") as ff:
        batch_mtree_entries_data = json.load(ff)
    if not os.path.exists(output_path):
        os.mkdir(output_path)

    for i in range(len(batch_mtree_entries_data)):
        user_entry_data = bytes.fromhex(batch_mtree_entries_data[str(i)])

        batch_inclusion_proof = {}
        batch_inclusion_proof["batch_mtree_root"] = batch_mtree[1]
        batch_inclusion_proof["user_id"] = str(
            user_entry_data[(coin_num+1)*32:(coin_num+2)*32].hex())
        batch_inclusion_proof["total_value"] = str(
            int.from_bytes(user_entry_data[:32], 'big'))
        j = 0
        for coin in coins:
            value = int.from_bytes(user_entry_data[(j+1)*32:(j+2)*32], 'big')
            if value > MAX_USER_VALUE:
                value = value - MODULUS
            batch_inclusion_proof[coin] = str(value)
            j += 1
        batch_inclusion_proof["random_number"] = str(
            user_entry_data[len(user_entry_data)-32:].hex())
        batch_inclusion_proof["user_index"] = i
        batch_inclusion_proof["batch_index"] = batch_index
        batch_inclusion_proof["uts"] = uts

        batch_inclusion_proof["merkle_path"] = mk_branch(
            batch_mtree, (uts * (i + 1) + uts-2) * EXTENSION_FACTOR)

        inclusion_proof = {
            "batch_inclusion_proof": batch_inclusion_proof,
            "trunk_inclusion_proof": trunk_inclusion_proof,
            "config": config_json
        }
        with open(output_path + "user_%d_inclusion_proof.json" % i, "w") as ff:
            json.dump(inclusion_proof, ff)

    # print("mk inclusion proof in %.4f sec" %(time.time() - start_time))

    return

# batch_index: batch index in trunk
# input_path: inclusion proof path, xxx/inclusion_proof_data/a1/


def verify_inclusion_proof(input_path):
    for root, dirs, files in os.walk(input_path):
        for proof_file in files:
            if re.search("inclusion_proof.json", proof_file):
                verify_single_inclusion_proof(input_path + proof_file)
    return


def verify_single_inclusion_proof(proof_file):
    with open(proof_file, "r") as ff:
        inclusion_proof = json.load(ff)

        coins = inclusion_proof["config"]["coins"]

        batch_inclusion_proof = inclusion_proof["batch_inclusion_proof"]
        batch_index = batch_inclusion_proof["batch_index"]
        user_index = batch_inclusion_proof["user_index"]
        uts = batch_inclusion_proof["uts"]
        user_leaf = verify_branch(bytes.fromhex(batch_inclusion_proof["batch_mtree_root"]), (uts * (
            user_index + 1) + uts-2) * EXTENSION_FACTOR, hex_array_to_bytes(batch_inclusion_proof["merkle_path"]))
        check_sum_value = int(batch_inclusion_proof["total_value"])
        user_entry = check_sum_value.to_bytes(32, 'big')
        j = 0
        temp = b''
        for coin in coins:
            value = int(batch_inclusion_proof[coin]) % MODULUS
            check_sum_value -= value
            temp = temp + value.to_bytes(32, 'big')
            j += 1
        user_entry = user_entry + keccak_256(temp) + bytes.fromhex(
            batch_inclusion_proof["user_id"]) + bytes.fromhex(batch_inclusion_proof["random_number"])
        assert check_sum_value == 0
        assert user_leaf == keccak_256(user_entry)

        trunk_inclusion_proof = inclusion_proof["trunk_inclusion_proof"]
        batch_leaf = verify_branch(bytes.fromhex(trunk_inclusion_proof["trunk_mtree_root"]), (UTS_FOR_TRUNK * (
            batch_index + 1) + UTS_FOR_TRUNK-2) * EXTENSION_FACTOR, hex_array_to_bytes(trunk_inclusion_proof["merkle_path"]))
        check_sum_value = int(trunk_inclusion_proof["total_value"])
        batch_entry = check_sum_value.to_bytes(32, 'big')
        j = 0
        temp = b''
        for coin in coins:
            value = int(trunk_inclusion_proof[coin]) % MODULUS
            check_sum_value -= value
            temp = temp + value.to_bytes(32, 'big')
            j += 1
        batch_entry = batch_entry + keccak_256(temp) + bytes.fromhex(
            trunk_inclusion_proof["batch_id"]) + bytes.fromhex(trunk_inclusion_proof["random_number"])
        assert check_sum_value == 0
        assert batch_leaf == keccak_256(batch_entry)

        assert trunk_inclusion_proof["batch_id"] == batch_inclusion_proof["batch_mtree_root"]

    return
