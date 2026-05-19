
import os
import pickle
from constants import PRECOMPUTED_PATH, MODULUS, NONRESIDUE, EXTENSION_FACTOR
from poly_utils import PrimeField
from utils import *
from fft import fft


def get_precomputed(user_num, uts):
    if not os.path.exists(PRECOMPUTED_PATH):
        os.makedirs(PRECOMPUTED_PATH)
    file_name = PRECOMPUTED_PATH + str(user_num) + "-" + str(uts) + ".json"
    if os.path.exists(file_name):
        with open(file_name, "rb") as f:
            data = pickle.load(f)
            return data
    else:
        data = precompute(file_name, user_num, uts)
        return data


def precompute(file_name, user_num, uts):
    f = PrimeField(MODULUS)
    steps = user_num * uts
    precision = steps * EXTENSION_FACTOR
    G2 = f.exp(NONRESIDUE, (MODULUS - 1) // precision)
    G1 = f.exp(G2, EXTENSION_FACTOR)
    domain = get_power_cycle(G1, MODULUS)
    xs = get_power_cycle(G2, MODULUS)
    domain_last_step_position = domain[steps - 1]

    z_eval = [domain[(i * user_num) % steps] - 1 for i in range(steps)]
    zerofier_poly1 = f.fit(fft(z_eval, MODULUS, G1, inv=True))

    z_num_eval = [(xs[(i * steps) % precision] - 1) %
                  MODULUS for i in range(precision)]
    z_den_eval = [(f.mul(f.sub(xs[(i * user_num) % precision], xs[(uts-2) * EXTENSION_FACTOR * user_num]),
                         f.sub(xs[(i * user_num) % precision], xs[(uts-1) * EXTENSION_FACTOR * user_num]))) for i in range(precision)]
    a = f.fit(fft(z_num_eval, MODULUS, G2, inv=True))
    b = f.fit(fft(z_den_eval, MODULUS, G2, inv=True))
    zerofier_poly2 = f.div_polys(a, b)

    z_num_eval = [f.sub(domain[(i * user_num) % steps], domain[(uts-1)
                        * user_num]) for i in range(steps)]
    z_den_eval = [f.sub(domain[i], domain_last_step_position)
                  for i in range(steps)]
    a = f.fit(fft(z_num_eval, MODULUS, G1, inv=True))
    b = f.fit(fft(z_den_eval, MODULUS, G1, inv=True))
    zerofier_poly3 = f.div_polys(a, b)

    zerofier_poly4 = f.mul_polys(
        [-domain[(uts-1)], 1], [-domain_last_step_position, 1])

    z_num_eval5 = [f.sub(domain[(i * user_num) % steps], domain[(uts-2)
                                                                * user_num]) for i in range(steps)]
    zerofier_poly5 = f.fit(fft(z_num_eval5, MODULUS, G1, inv=True))

    data = (G1, G2, domain, zerofier_poly1, zerofier_poly2,
            zerofier_poly3, zerofier_poly4, zerofier_poly5)
    with open(file_name, "wb") as f:
        pickle.dump(data, f)

    return data
