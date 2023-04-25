from keccak_merkle_tree import merkelize as _merkelize
from keccak_merkle_tree import mk_branch as _mk_branch
from keccak_merkle_tree import verify_branch as _verify_branch
from keccak_merkle_tree import mk_multi_branch as _mk_multi_branch
from keccak_merkle_tree import verify_multi_branch as _verify_multi_branch
from keccak_merkle_tree import keccak_256

def permute4_values(values):
    o = []
    ld4 = len(values) // 4
    for i in range(ld4):
        o.extend([values[i], values[i + ld4], values[i + ld4 * 2], values[i + ld4 * 3]])
    return o

def permute4_index(x, L):
    ld4 = L // 4
    return x//ld4 + 4 * (x % ld4)

def permute4_indices(xs, L):
    ld4 = L // 4
    return [x//ld4 + 4 * (x % ld4) for x in xs]

def merkelize(L):
    return _merkelize(permute4_values(L))

def mk_branch(tree, index):
    return _mk_branch(tree, permute4_index(index, len(tree) // 2))

def verify_branch(root, index, proof, output_as_int=False):
    return _verify_branch(root, permute4_index(index, 2**len(proof) // 2), proof, output_as_int)

def mk_multi_branch(tree, indices):
    return _mk_multi_branch(tree, permute4_indices(indices, len(tree) // 2))    

def verify_multi_branch(root, indices, proof):
    return _verify_multi_branch(root, permute4_indices(indices, 2**len(proof[0]) // 2), proof)



# L = [0,1,2,3,4,5,6,7]
# mtree = merkelize(L)
# m = []
# for i in range(len(mtree)):
#     m.append(mtree[i].hex())
    
# print("mtree", m)    
# print("root", mtree[1].hex())

# proof = mk_branch(mtree, 2)
# p = []
# for i in range(len(proof)):
#     p.append(proof[i].hex())
# print("proof", p)
# print("proof.length", len(proof))

# leaf = verify_branch(mtree[1], 2, proof)
# print("leaf", leaf)