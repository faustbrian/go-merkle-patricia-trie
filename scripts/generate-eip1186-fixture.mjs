// Future deterministic producer; this does not identify the historical writer.
// Install the unchanged authenticated package lock with npm ci --ignore-scripts.
// Emit only the preserved, checksummed corpus; redirect stdout to an owned
// temporary file and compare it before any explicitly reviewed corpus update.
import { createHash } from "node:crypto";
import {
  MerklePatriciaTrie,
  createMerkleProof,
  verifyMPTWithMerkleProof,
} from "@ethereumjs/mpt";
import { RLP } from "@ethereumjs/rlp";

if (process.argv.length !== 2) {
  throw new Error("The fixed fixture producer accepts no arguments");
}

const hex = (value) => Buffer.from(value).toString("hex");
const word = (size, index, value) => {
  const result = new Uint8Array(size);
  result[index] = value;
  return result;
};
const codeHash = Buffer.from(
  "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
  "hex",
);
const emptyRoot = new MerklePatriciaTrie().root();
const slot = word(32, 31, 7);
const absentSlot = word(32, 31, 8);
const storageValue = new Uint8Array([1, 0x80]);
const storage = new MerklePatriciaTrie({ useKeyHashing: true });
await storage.put(slot, RLP.encode(storageValue));
await storage.put(word(32, 31, 9), RLP.encode(new Uint8Array([1])));

const address = word(20, 19, 0xaa);
const absentAddress = word(20, 0, 0x99);
const account = RLP.encode([
  new Uint8Array([3]),
  storageValue,
  storage.root(),
  codeHash,
]);
const state = new MerklePatriciaTrie({ useKeyHashing: true });
await state.put(address, account);
await state.put(
  word(20, 0, 0x10),
  RLP.encode([new Uint8Array(), new Uint8Array(), emptyRoot, codeHash]),
);

async function proof(trie, key, expectedValue) {
  const nodes = await createMerkleProof(trie, key);
  const value = await verifyMPTWithMerkleProof(trie, trie.root(), key, nodes);
  if (
    expectedValue === null
      ? value !== null
      : value === null || hex(value) !== hex(expectedValue)
  ) {
    throw new Error("Locked MPT oracle rejected a fixed fixture claim");
  }
  return nodes.map(hex);
}

const fixture = {
  absentAccountProof: await proof(state, absentAddress, null),
  absentAddress: hex(absentAddress),
  absentSlot: hex(absentSlot),
  absentStorageProof: await proof(storage, absentSlot, null),
  account: hex(account),
  accountProof: await proof(state, address, account),
  address: hex(address),
  format: 1,
  slot: hex(slot),
  stateRoot: hex(state.root()),
  storageProof: await proof(storage, slot, RLP.encode(storageValue)),
  storageRoot: hex(storage.root()),
  storageValue: hex(storageValue),
};
const serialized = Buffer.from(`${JSON.stringify(fixture, null, 2)}\n`);
const digest = createHash("sha256").update(serialized).digest("hex");
if (
  serialized.length !== 1769 ||
  digest !== "9ef6b4b3a4b740172a6372c717277de725ace3d8ac0af958d4087bcc13f6d000"
) {
  throw new Error("Produced corpus differs from the preserved fixture");
}
process.stdout.write(serialized);
