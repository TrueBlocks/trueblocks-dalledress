export type Address = string;

export const isValidAddress = (address: string): address is Address => {
  return /^0x[a-fA-F0-9]{40}$/.test(address) || address.endsWith('.eth');
};

export const createAddress = (address: string): Address => {
  if (!isValidAddress(address)) {
    throw new Error(`Invalid address format: ${address}`);
  }
  return address as Address;
};
