# This demo is based on Vitalik's proof of solvency proposal and implemented using the STARK proof system.
# See https://vitalik.ca/general/2022/11/19/proof_of_solvency.html
# It provides users with proofs that constrain the sum of all assets and the non-negativity of their net asset value.

from permuted_tree import merkelize, hash, mk_multi_branch, verify_multi_branch, mk_branch, verify_branch
from poly_utils import PrimeField
from fft import fft
from fri import prove_low_degree, verify_low_degree_proof
from utils import *
from constants import *
import time
import gc

f = PrimeField(MODULUS)

# ids: array of user ids
# values: array of user values of all coins
# uts: user trace size, number of data rows for each user, for non-negative proof
# data_path: user data path


def mk_por_proof(ids, values, uts, data_path, main_coins_num, coins):
    start_time = time.time()

    # tranform data into the field, -1 => MODULUS -1
    transform_into_field(values, MODULUS)
    assert is_a_power_of_2(uts) and uts <= MAX_UTS, "invalid uts"

    if not is_a_power_of_2(len(ids)+1):
        ids, values = pad(ids, values, MAX_USER_NUM_FOR_ONE_BATCH)
    user_num = len(ids)+1
    steps = uts * user_num
    precision = steps * EXTENSION_FACTOR

    ids, values = extend_user_data(ids, values, uts)
    sum_trace, sum_values, values = get_sum_trace(values, uts, MODULUS)

    # get generators
    G2 = f.exp(NONRESIDUE, (MODULUS - 1) // precision)
    skips = precision // steps
    G1 = f.exp(G2, skips)

    # origin domain before extension
    domain = get_power_cycle(G1, MODULUS)
    # get x coordinates for fft
    xs = get_power_cycle(G2, MODULUS)
    last_step_position = xs[(steps - 1) * EXTENSION_FACTOR]

    # get poly and eval for sum_trace
    t_poly = fft(sum_trace, MODULUS, G1, inv=True)
    t_eval = fft(t_poly, MODULUS, G2)
    # del t_poly

    # get poly and eval for values of each coin
    b_eval = []
    b_poly = []
    for i in range(len(values)):
        b_poly.append(fft(values[i], MODULUS, G1, inv=True))
        b_eval.append(fft(b_poly[-1], MODULUS, G2))
    # del b_poly
    # gc.collect()

    # trace constraints evaluations
    tc_eval = []
    # constraint 1: The first row of each user's trace should be 0, to ensure the non-negativity of user's net assets value
    # t[uts*EXTENSION_FACTOR*i] = 0, 0<=i<=user_num-1,
    # z1(x) = (x-xs[uts*EXTENSION_FACTOR*0])(x-xs[uts*EXTENSION_FACTOR*1])...(x-xs[uts*EXTENSION_FACTOR*(user_num-1)])
    # z1(x) = x^user_num - 1
    z_num_eval = [xs[(i * user_num) % precision] - 1 for i in range(precision)]
    z_eval = f.multi_inv(z_num_eval)
    tc_eval.append([tp * z % MODULUS for tp, z in zip(t_eval, z_eval)])

    zerofier1 = fft(z_num_eval, MODULUS, G2, inv=True)

    # constraint 2:  row_value = next_row_value // 4, in case user's net assets value less than 4^(uts-2), the first row of each user's trace should be 0
    # (t[i + EXTENSION_FACTOR] - 4*t[i])*(t[i + EXTENSION_FACTOR] - 4*t[i] - 1)*(t[i + EXTENSION_FACTOR] - 4*t[i] - 2)*(t[i + EXTENSION_FACTOR] - 4*t[i] - 3) = 0,
    # 0<=i<=steps-1 && (i mod uts*EXTENSION_FACTOR != {(uts-2)*EXTENSION_FACTOR,(uts-1)*EXTENSION_FACTOR}, i in range(precision)
    # z2(x) = (x^steps -1)/((x^user_num -  G2^((uts-2)*EXTENSION_FACTOR*user_num))(x^user_num - G2^((uts-1)*EXTENSION_FACTOR*user_num)))
    c_num_eval = [f.mul(f.mul(f.sub(t_eval[(i + EXTENSION_FACTOR) % precision], 4 * t_eval[i]),
                              (f.sub(t_eval[(i + EXTENSION_FACTOR) % precision], 4 * t_eval[i]) - 1)),
                        (f.mul((f.sub(t_eval[(i + EXTENSION_FACTOR) % precision], 4 * t_eval[i]) - 2),
                               (f.sub(t_eval[(i + EXTENSION_FACTOR) % precision], 4 * t_eval[i]) - 3)))) for i in range(precision)]
    z_num_eval = [(xs[(i * steps) % precision] - 1) %
                  MODULUS for i in range(precision)]
    z_num_inv = f.multi_inv(z_num_eval)
    z_den_eval = [(f.mul(f.sub(xs[(i * user_num) % precision], xs[(uts-2) * EXTENSION_FACTOR * user_num]),
                         f.sub(xs[(i * user_num) % precision], xs[(uts-1) * EXTENSION_FACTOR * user_num]))) for i in range(precision)]
    tc_eval.append([f.mul(f.mul(cn, zi), zd)
                   for cn, zi, zd in zip(c_num_eval, z_num_inv, z_den_eval)])

    # constraint 3: User's net asset value accumulation
    # t(i + uts*EXTENSION_FACTOR) = t(i + (uts-1)*EXTENSION_FACTOR) + t(i), i mod uts*EXTENSION_FACTOR == (uts-1)*EXTENSION_FACTOR， and i != last_step_position，
    # z3(x) = (x^user_num - G2^((uts-1) * EXTENSION_FACTOR * user_num))/(x - last_step_position)
    c_num_eval = [f.sub(f.sub(t_eval[(i + uts*EXTENSION_FACTOR) % precision], t_eval[(
        i + (uts-1)*EXTENSION_FACTOR) % precision]), t_eval[i]) for i in range(precision)]
    z_num_eval = [f.sub(xs[(i * user_num) % precision], xs[(uts-1)
                        * EXTENSION_FACTOR * user_num]) for i in range(precision)]
    z_num_inv = f.multi_inv(z_num_eval)
    z_den_eval = [f.sub(xs[i], last_step_position) for i in range(precision)]
    tc_eval.append([f.mul(f.mul(cn, zi), zd)
                   for cn, zi, zd in zip(c_num_eval, z_num_inv, z_den_eval)])

    # constraint 4: The initial accumulation should be 0, the last accumulation should total value of all assets of all users
    # t((uts-1)*EXTENSION_FACTOR) = 0， t(last_step_position) = sum_values[-1]
    # z4(x) = (x-xs[(uts-1)*EXTENSION_FACTOR])(x-last_step_position)
    interpolant = f.lagrange_interp_2(
        [xs[(uts-1)*EXTENSION_FACTOR], last_step_position], [0, sum_values[-1]])
    i_eval = [f.eval_poly_at(interpolant, x) for x in xs]
    z_poly = f.mul_polys([-xs[(uts-1)*EXTENSION_FACTOR],
                         1], [-last_step_position, 1])
    z_eval = f.multi_inv([f.eval_poly_at(z_poly, x) for x in xs])
    tc_eval.append([f.mul(f.sub(t, i), z) % MODULUS for t,
                   i, z in zip(t_eval, i_eval, z_eval)])

    # constraint 5: user's sum values of each coin should be the t_eval
    #  t(i) - sum(b_eval[i][j] for j in range(len(b_eval))) = 0, i mod uts*EXTENSION_FACTOR == (uts-2)*EXTENSION_FACTOR
    # z5(x) = (x^user_num - G2^((uts-2) * EXTENSION_FACTOR * user_num))
    c_num_eval = [f.sub(t_eval[j], sum(b_eval[i][j]
                                       for i in range(len(b_eval)))) for j in range(precision)]

    z_num_eval5 = [f.sub(xs[(i * user_num) % precision], xs[(uts-2)
                                                            * EXTENSION_FACTOR * user_num]) for i in range(precision)]
    z_num_inv = f.multi_inv(z_num_eval5)
    tc_eval.append([f.mul(cn, zi) for cn, zi in zip(c_num_eval, z_num_inv)])

    # values constraints eval
    cc_eval = []
    for i in range(main_coins_num):
        # cc_constraint 1: Accumulation of each coin
        # b_eval[i][j + uts*EXTENSION_FACTOR] = b_eval[i][j + (uts-1)*EXTENSION_FACTOR] + b_eval[i][j], j mod uts*EXTENSION_FACTOR == (uts-1)*EXTENSION_FACTOR，and j != last_step_position，
        # z(x) = (x^user_num - G2^((uts-1) * EXTENSION_FACTOR * user_num))/(x - last_step_position)
        c_num_eval = [f.sub(f.sub(b_eval[i][(j + uts*EXTENSION_FACTOR) % precision], b_eval[i][(
            j + (uts-1)*EXTENSION_FACTOR) % precision]), b_eval[i][j]) for j in range(precision)]
        z_num_eval = [f.sub(xs[(i * user_num) % precision], xs[(uts-1)
                            * EXTENSION_FACTOR * user_num]) for i in range(precision)]
        z_num_inv = f.multi_inv(z_num_eval)
        z_den_eval = [f.sub(xs[i], last_step_position)
                      for i in range(precision)]
        cc_eval.append([f.mul(f.mul(cn, zi), zd) % MODULUS for cn,
                       zi, zd in zip(c_num_eval, z_num_inv, z_den_eval)])

        # cc_constraints 2: The initial accumulation should be 0, the last accumulation should total value of this coin of all users
        # b_eval[i](uts-1)*EXTENSION_FACTOR) = 0, b_eval[i](last_step_position) = sum_amount[i]
        # z(x) = (x-xs[(uts-1)*EXTENSION_FACTOR])(x-last_step_position)
        interpolant = f.lagrange_interp_2(
            [xs[(uts-1)*EXTENSION_FACTOR], last_step_position], [0, sum_values[i]])
        i_eval = [f.eval_poly_at(interpolant, x) for x in xs]
        cc_eval.append([f.mul(f.sub(t, i), zi)
                       for t, i, zi in zip(b_eval[i], i_eval, z_eval)])

    del i_eval, interpolant, z_eval, c_num_eval, z_poly, z_den_eval, z_num_inv, z_num_eval
    gc.collect()

    user_random = [int.from_bytes(hash(r), 'big')
                   for r in get_entries([tc_eval, cc_eval])]
    id_eval = sum([[x] + [0]*(EXTENSION_FACTOR-1) for x in ids], [])

    user_entry_data = []
    for i in range(1, user_num):
        index = (i * uts + uts - 2) * EXTENSION_FACTOR
        data = t_eval[index].to_bytes(32, 'big') + get_entry_data(b_eval, index) + \
            id_eval[index].to_bytes(32, 'big') + \
            user_random[index].to_bytes(32, 'big')
        user_entry_data.append(data)
    save_mtree_entries_data(data_path, user_entry_data)
    del user_entry_data
    gc.collect()

    mtree = merkelize(get_leaves([t_eval, b_eval, id_eval, user_random]))

    pow_nonce = proof_of_work(mtree[1], POW_BITS)

    # linearly combination
    G2_to_the_steps = f.exp(G2, 3 * steps)
    powers = [1]
    for i in range(1, precision):
        powers.append(powers[-1] * G2_to_the_steps % MODULUS)

    l_eval = calculate_l(pow_nonce, powers, [
                         t_eval, b_eval, id_eval, tc_eval[0], tc_eval[2:], cc_eval], MODULUS)
    l_eval = [(l + c) % MODULUS for l, c in zip(l_eval, tc_eval[1])]
    l_mtree = merkelize(l_eval)

    # sampling
    positions = get_pseudorandom_indices(l_mtree[1], precision, SPOT_CHECK_SECURITY_FACTOR,
                                         exclude_multiples_of=EXTENSION_FACTOR)
    aug_positions = sum([[x, (x + skips) % precision, (x + (uts-1)*skips) %
                        precision, (x + uts*skips) % precision] for x in positions], [])

    sampled_entries_data = []
    for i in range(SPOT_CHECK_SECURITY_FACTOR):
        sampled_entries_data = sampled_entries_data + [get_entry_data([t_eval, b_eval, id_eval, tc_eval, cc_eval], aug_positions[4*i]),
                                                       [t_eval[aug_positions[4*i+1]].to_bytes(32, 'big'),
                                                        hash(get_entry_data(
                                                            b_eval, aug_positions[4*i+1])),
                                                        id_eval[aug_positions[4*i+1]
                                                                ].to_bytes(32, 'big'),
                                                        user_random[aug_positions[4*i+1]].to_bytes(32, 'big')],
                                                       [t_eval[aug_positions[4*i+2]].to_bytes(32, 'big'),
                                                        get_entry_data(
                                                            b_eval, aug_positions[4*i+2]),
                                                        id_eval[aug_positions[4*i+2]
                                                                ].to_bytes(32, 'big'),
                                                        user_random[aug_positions[4*i+2]].to_bytes(32, 'big')],
                                                       [t_eval[aug_positions[4*i+3]].to_bytes(32, 'big'),
                                                        get_entry_data(
                                                            b_eval, aug_positions[4*i+3]),
                                                        id_eval[aug_positions[4*i+3]
                                                                ].to_bytes(32, 'big'),
                                                        user_random[aug_positions[4*i+3]].to_bytes(32, 'big')]]
    del t_eval, b_eval, id_eval, tc_eval, cc_eval
    gc.collect()
    # return the merkle roots, the spot check merkle proofs, and low-degree proofs
    sum_proof = [steps,
                 uts,
                 mtree[1],
                 l_mtree[1],
                 pow_nonce,
                 mk_multi_branch(mtree, aug_positions),
                 sampled_entries_data,
                 mk_multi_branch(l_mtree, positions),
                 prove_low_degree(l_eval, G2, 4*steps, MODULUS, exclude_multiples_of=EXTENSION_FACTOR)]

    save_data(data_path, sum_proof, mtree, sum_values, coins)
    # print("mk por proof in %.4f sec: " % (time.time() - start_time))
    return

# sum_values: The sum values of each coin ant total coin that prover claimed
# proof: The proof for the sum amounts


def verify_por_proof(sum_values, proof, main_coins_num):
    start_time = time.time()
    check_sum_values(sum_values, MODULUS)
    coins_num = len(sum_values) - 1
    steps, uts, m_root, l_root, pow_nonce, main_branches, mtree_entries_data, linear_comb_branches, fri_proof = proof
    assert steps <= 2**32 // EXTENSION_FACTOR, "invalid steps: too large"
    assert is_a_power_of_2(steps), "invalid steps: should be a power of 2"
    assert int.from_bytes(hash(m_root + pow_nonce),
                          'big') >> (256 - POW_BITS) == 0, "invalid proof of work"

    precision = steps * EXTENSION_FACTOR
    user_num = steps // uts

    G2 = f.exp(NONRESIDUE, (MODULUS-1)//precision)
    skips = precision // steps
    G1 = f.exp(G2, skips)
    assert verify_low_degree_proof(
        l_root, G2, fri_proof, 4*steps, MODULUS, exclude_multiples_of=EXTENSION_FACTOR)

    # performs the spot checks
    k = [int.from_bytes(hash(pow_nonce + i.to_bytes(32, 'big')),
                        'big') % MODULUS for i in range(6 * coins_num + 12)]

    positions = get_pseudorandom_indices(l_root, precision, SPOT_CHECK_SECURITY_FACTOR,
                                         exclude_multiples_of=EXTENSION_FACTOR)
    aug_positions = sum([[x, (x + skips) % precision, (x + (uts-1)*skips) %
                        precision, (x + uts*skips) % precision] for x in positions], [])
    last_step_position = f.exp(G2, (steps - 1) * skips)

    main_branch_leaves = verify_multi_branch(
        m_root, aug_positions, main_branches)
    check_entry_hash(main_branch_leaves, mtree_entries_data,
                     main_coins_num, MODULUS)

    linear_comb_branch_leaves = verify_multi_branch(
        l_root, positions, linear_comb_branches)

    for i, pos in enumerate(positions):
        x = f.exp(G2, pos)
        x_to_the_steps = f.exp(x, 3*steps)
        mbranch1 = mtree_entries_data[i*4]
        mbranch2 = mtree_entries_data[i*4+1]
        mbranch3 = mtree_entries_data[i*4+2]
        mbranch4 = mtree_entries_data[i*4+3]

        l_of_x = int.from_bytes(linear_comb_branch_leaves[i], 'big')

        t_of_x = int.from_bytes(mbranch1[:32], 'big')
        b_of_x = [int.from_bytes(mbranch1[32+32*i:64+32*i], 'big')
                  for i in range(coins_num)]
        id_of_x = int.from_bytes(
            mbranch1[32+32*coins_num:64+32*coins_num], 'big')
        tc_of_x = [int.from_bytes(
            mbranch1[64+32*(coins_num+i):96+32*(coins_num+i)], 'big') for i in range(5)]
        cc_of_x = [int.from_bytes(
            mbranch1[224+32*(coins_num+i):256+32*(coins_num+i)], 'big') for i in range(2*coins_num)]

        t_of_skips_x = int.from_bytes(mbranch2[0], 'big')
        t_of_uts_sub_1_skips_x = int.from_bytes(mbranch3[0], 'big')
        t_of_uts_skips_x = int.from_bytes(mbranch4[0], 'big')

        b_of_uts_sub_1_skips_x = [int.from_bytes(
            mbranch3[1][32*i:32+32*i], 'big') for i in range(coins_num)]
        b_of_uts_skips_x = [int.from_bytes(
            mbranch4[1][32*i:32+32*i], 'big') for i in range(coins_num)]

        # check constraint 1: t[uts*EXTENSION_FACTOR*i] = 0, 0<=i<=user_num-1
        # t[i] = c1[i] * z1[i]
        z1 = f.exp(x, user_num) - 1
        assert t_of_x == f.mul(tc_of_x[0], z1)

        # check constraint 2:  (t[i + EXTENSION_FACTOR] - 4*t[i])*(t[i + EXTENSION_FACTOR] - 4*t[i] - 1)*(t[i + EXTENSION_FACTOR] - 4*t[i] - 2)*(t[i + EXTENSION_FACTOR] - 4*t[i] - 3) = 0,
        # 0<=i<=steps-1 && (i mod uts*EXTENSION_FACTOR != {(uts-2)*EXTENSION_FACTOR,(uts-1)*EXTENSION_FACTOR}, i in range(precision)
        # (t[i + EXTENSION_FACTOR] - 4*t[i])*(t[i + EXTENSION_FACTOR] - 4*t[i] - 1)*(t[i + EXTENSION_FACTOR] - 4*t[i] - 2)*(t[i + EXTENSION_FACTOR] - 4*t[i] - 3) = c2[i] * z2[i]
        z2 = f.div(f.exp(x, steps) - 1,
                   f.mul(f.sub(f.exp(x, user_num), f.exp(G2, (uts-2) * EXTENSION_FACTOR * user_num)),
                         f.sub(f.exp(x, user_num), f.exp(G2, (uts-1) * EXTENSION_FACTOR * user_num))))
        assert f.mul(f.mul(f.sub(t_of_skips_x, 4 * t_of_x), f.sub(t_of_skips_x, 4 * t_of_x + 1)),
                     f.mul(f.sub(t_of_skips_x, 4 * t_of_x + 2), f.sub(t_of_skips_x, 4 * t_of_x + 3))) == f.mul(tc_of_x[1], z2)

        # check constraint 3: t(i + uts*EXTENSION_FACTOR) = t(i + (uts-1)*EXTENSION_FACTOR) + t(i),i mod uts*EXTENSION_FACTOR == (uts-1)*EXTENSION_FACTOR，i != last_step_position
        # t(i + uts*EXTENSION_FACTOR) - t(i + (uts-1)*EXTENSION_FACTOR) - t(i) = c3[i] * z3
        z3 = f.div(f.sub(f.exp(x, user_num), f.exp(G2, (uts-1) * EXTENSION_FACTOR * user_num)),
                   f.sub(x, last_step_position))
        assert f.sub(f.sub(t_of_uts_skips_x, t_of_uts_sub_1_skips_x),
                     t_of_x) == f.mul(tc_of_x[2], z3)

        # check constraint 4: t((uts-1)*EXTENSION_FACTOR) = 0， t(last_step_position) = sum_values[-1]
        # t[i] - interpolant_value = c4[i] * z4[i]
        interpolant = f.lagrange_interp_2(
            [f.exp(G2, (uts-1)*EXTENSION_FACTOR), last_step_position], [0, sum_values[-1]])
        interpolant_value = f.eval_poly_at(interpolant, x)
        z4_poly = f.mul_polys(
            [-f.exp(G2, ((uts-1)*EXTENSION_FACTOR)), 1], [-last_step_position, 1])
        z4 = f.eval_poly_at(z4_poly, x)
        assert f.sub(t_of_x, interpolant_value) == f.mul(tc_of_x[3], z4)

        # check constraint 5:t(i) - sum(b_eval[i][j] for j in range(len(b_eval))) = 0, i mod uts*EXTENSION_FACTOR == (uts-2)*EXTENSION_FACTOR
        # z5(x) = (x^user_num - G2^((uts-2) * EXTENSION_FACTOR * user_num))

        z5 = f.sub(f.exp(x, user_num), f.exp(
            G2, (uts-2) * EXTENSION_FACTOR * user_num))
        assert f.sub(t_of_x, sum(b_of_x)) == f.mul(tc_of_x[4], z5)

        # check coins constraint:
        for i in range(main_coins_num):
            # check cc_constraint 1: b_eval[i][j + uts*EXTENSION_FACTOR] = b_eval[i][j + (uts-1)*EXTENSION_FACTOR] + b_eval[i][j]
            # j mod uts*EXTENSION_FACTOR == (uts-1)*EXTENSION_FACTOR，and j != last_step_position
            assert f.sub(f.sub(b_of_uts_skips_x[i], b_of_uts_sub_1_skips_x[i]), b_of_x[i]) == f.mul(
                cc_of_x[2*i], z3)

            # check cc_constraints 2: b_eval[i](uts-1)*EXTENSION_FACTOR) = 0, b_eval[i](last_step_position) = 0
            interpolant = f.lagrange_interp_2(
                [f.exp(G2, (uts-1)*EXTENSION_FACTOR), last_step_position], [0, sum_values[i]])
            interpolant_value = f.eval_poly_at(interpolant, x)
            assert f.sub(b_of_x[i], interpolant_value) == f.mul(
                cc_of_x[2*i+1], z4)

        # check correctness of the linear combination
        assert verify_l(k, x_to_the_steps, l_of_x, [
                        t_of_x, b_of_x, id_of_x, tc_of_x[0], tc_of_x[2:], cc_of_x], tc_of_x[1], MODULUS)

    # print('Verified %d consistency checks' % SPOT_CHECK_SECURITY_FACTOR)
    # print('Verified sum proof in %.4f sec' % (time.time() - start_time))
    return True
