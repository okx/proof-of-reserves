from mk_and_verify_proofs import *
from utils import *
from constants import *
from zk_STARK_Validator import por_user_verify_proofs
from por_verify_tool import verify_sum_proofs, verify_all_inclusion_proof
import json
import time


def test_full_process():
    with open(CONFIG_PATH, "r") as ff:
        config_json = json.load(ff)
    init_user_data(15, 0, CONFIG_PATH)
    init_user_data(15, 1, CONFIG_PATH)
    init_user_data(15, 2, CONFIG_PATH)
    init_user_data(15, 3, CONFIG_PATH)

    mk_batch_proof(16, "./user_data/batch0.json",
                   "./sum_proof_data/batches/a0/", CONFIG_PATH)
    mk_batch_proof(16, "./user_data/batch1.json",
                   "./sum_proof_data/batches/a1/", CONFIG_PATH)
    mk_batch_proof(16, "./user_data/batch2.json",
                   "./sum_proof_data/batches/a2/", CONFIG_PATH)
    mk_batch_proof(32, "./user_data/batch3.json",
                   "./sum_proof_data/batches/b0/", CONFIG_PATH)

    mk_trunk_proof("./sum_proof_data/batches/",
                   "./sum_proof_data/trunk/", CONFIG_PATH)

    verify_batch_proof("./sum_proof_data/batches/a0/", config_json)
    verify_batch_proof("./sum_proof_data/batches/a1/", config_json)
    verify_batch_proof("./sum_proof_data/batches/a2/", config_json)
    verify_batch_proof("./sum_proof_data/batches/b0/", config_json)

    verify_trunk_proof("./sum_proof_data/trunk/", config_json)

    mk_inclusion_proof(0, 16, "./sum_proof_data/batches/a0/",
                       "./sum_proof_data/trunk/", "./inclusion_proof_data/a0/", CONFIG_PATH)
    mk_inclusion_proof(1, 16, "./sum_proof_data/batches/a1/",
                       "./sum_proof_data/trunk/", "./inclusion_proof_data/a1/", CONFIG_PATH)
    mk_inclusion_proof(2, 16, "./sum_proof_data/batches/a2/",
                       "./sum_proof_data/trunk/", "./inclusion_proof_data/a2/", CONFIG_PATH)
    mk_inclusion_proof(3, 32, "./sum_proof_data/batches/b0/",
                       "./sum_proof_data/trunk/", "./inclusion_proof_data/b0/", CONFIG_PATH)

    verify_inclusion_proof("./inclusion_proof_data/a0/")
    verify_inclusion_proof("./inclusion_proof_data/a1/")
    verify_inclusion_proof("./inclusion_proof_data/a2/")
    verify_inclusion_proof("./inclusion_proof_data/b0/")

    # por_user_verify_proofs()

    print("Full Process Test Passed")

    return


def test_try_invalid_sum_value():
    with open("./sum_proof_data/batches/a0/sum_values.json", "r") as ff:
        sum_values_json = json.load(ff)
        sum_values_json["total_value"] += 1

    with open("./sum_proof_data/batches/a0/sum_values.json", "w") as ff:
        json.dump(sum_values_json, ff)

    with open(CONFIG_PATH, "r") as ff:
        config_json = json.load(ff)
    try:
        verify_batch_proof("./sum_proof_data/batches/a0/", config_json)
    except:
        print("Invalid Sum Value")

    return


def test_negative_value_with_positive_net_value():
    with open("./user_data/batch0.json", "r") as ff:
        user_data_json = json.load(ff)
        user_data_json[0]["BTC"] = str(-1)

    with open("./user_data/batch0.json", "w") as ff:
        json.dump(user_data_json, ff)

    with open(CONFIG_PATH, "r") as ff:
        config_json = json.load(ff)

    mk_batch_proof(16, "./user_data/batch0.json",
                   "./sum_proof_data/batches/a0/", CONFIG_PATH)

    verify_batch_proof("./sum_proof_data/batches/a0/", config_json)

    print("Negative Value with Positive Net Value Is Allowed")

    init_user_data(15, 0, CONFIG_PATH)

    return


def test_negative_net_value():
    with open(CONFIG_PATH, "r") as ff:
        config_json = json.load(ff)
        coins = config_json["coins"]

    with open("./user_data/batch0.json", "r") as ff:
        user_data_json = json.load(ff)
        data_copy = user_data_json
        for coin in coins:
            user_data_json[0][coin] = str(- int(user_data_json[0][coin]))

    with open("./user_data/batch0.json", "w") as ff:
        json.dump(user_data_json, ff)

    mk_batch_proof(16, "./user_data/batch0.json",
                   "./sum_proof_data/batches/a0/", CONFIG_PATH)

    try:
        verify_batch_proof("./sum_proof_data/batches/a0/", config_json)
    except:
        print("Invalid Net Value")

    init_user_data(15, 0, CONFIG_PATH)

    return


def test_invalid_inclusion_proof():
    with open("./inclusion_proof_data/a0/user_0_inclusion_proof.json", "r") as ff:
        inclusion_proof_json = json.load(ff)
        inclusion_proof_json["batch_inclusion_proof"]["total_value"] = str(
            int(inclusion_proof_json["batch_inclusion_proof"]["total_value"]) + 1)

    with open("./inclusion_proof_data/a0/user_0_inclusion_proof.json", "w") as ff:
        json.dump(inclusion_proof_json, ff)

    try:
        verify_inclusion_proof("./inclusion_proof_data/a0/")
    except:
        print("Invalid Inclusion Proof")

    mk_inclusion_proof(0, 16, "./sum_proof_data/batches/a0/",
                       "./sum_proof_data/trunk/", "./inclusion_proof_data/a0/", CONFIG_PATH)

    return


if __name__ == '__main__':
    time0 = time.time()
    test_full_process()
    test_try_invalid_sum_value()
    test_negative_value_with_positive_net_value()
    test_negative_net_value()
    test_invalid_inclusion_proof()
    test_full_process()
    verify_sum_proofs(CONFIG_PATH)
    verify_all_inclusion_proof()
    print("all test finished in %d sec" % (time.time()-time0))
