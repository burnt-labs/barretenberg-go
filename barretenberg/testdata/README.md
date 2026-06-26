# Test Vectors

Binary test vectors for UltraHonk proof verification, generated with Aztec Barretenberg v5.0.0-rc.1.

## Files

- `statics/vk` — Verification key (binary, UltraHonk format)
- `statics/proof` — Proof (binary)
- `statics/public_inputs` — Concatenated 32-byte field elements (big-endian)

## Regenerating

Requires [Noir](https://noir-lang.org/) (nargo) and [Barretenberg](https://github.com/AztecProtocol/aztec-packages) CLI (bb) v5.0.0-rc.1.

1. Install nargo and bb CLI at version v5.0.0-rc.1
2. Create a simple Noir circuit (e.g., `x * x == y`)
3. Generate proof and verification key:

```bash
nargo compile
bb write_vk -b target/circuit.json -o testdata/statics/vk
nargo prove
bb prove -b target/circuit.json -w target/witness.gz -o testdata/statics/proof
# Extract public inputs from the witness
```

See the original generation script for details:
https://github.com/burnt-labs/xion (branch: feature/barrentenberg-go-bindings, file: x/zk/barretenberg/testdata/generate.sh)
