import { types } from '@models';

export interface TransactionData {
  to: string;
  function: types.Function;
  inputs: TransactionInput[];
  value?: string; // ETH value for payable functions
}

export interface TransactionInput {
  name: string;
  type: string;
  value: string;
}

export interface PreparedTransaction {
  to: string;
  data: string;
  value: string;
  gas?: string;
  gasPrice?: string;
}

/**
 * Encodes function parameters for contract calls
 */
export const encodeParameters = (
  functionAbi: types.Function,
  inputs: TransactionInput[],
): string => {
  // TODO: Implement actual ABI encoding
  // This is a placeholder that would use ethers.js or web3.js for encoding

  const encodedParams = inputs
    .map((input) => {
      return encodeParameter(input.type, input.value);
    })
    .join('');

  // Function selector (first 4 bytes of keccak256 hash of function signature)
  const _functionSignature = `${functionAbi.name}(${functionAbi.inputs
    .map((i) => i.type)
    .join(',')})`;
  // TODO: Calculate actual keccak256 hash
  const selector = '0x12345678'; // Placeholder

  return selector + encodedParams;
};

/**
 * Encodes a single parameter value
 */
export const encodeParameter = (type: string, value: string): string => {
  // TODO: Implement actual parameter encoding based on type
  // This is a placeholder implementation

  if (type === 'address') {
    return value.replace('0x', '').padStart(64, '0');
  }

  if (type.startsWith('uint') || type.startsWith('int')) {
    const numValue = BigInt(value);
    return numValue.toString(16).padStart(64, '0');
  }

  if (type === 'bool') {
    return value.toLowerCase() === 'true'
      ? '1'.padStart(64, '0')
      : '0'.padStart(64, '0');
  }

  if (type === 'string') {
    // TODO: Implement proper string encoding
    return new TextEncoder()
      .encode(value)
      .toString()
      .replace(/[^0-9a-fA-F]/g, '')
      .padStart(64, '0');
  }

  // For arrays and complex types, this would need more sophisticated encoding
  return value.replace(/[^0-9a-fA-F]/g, '').padStart(64, '0');
};

/**
 * Builds a transaction from form data
 */
export const buildTransaction = (
  contractAddress: string,
  functionAbi: types.Function,
  inputs: TransactionInput[],
  ethValue?: string,
): TransactionData => {
  return {
    to: contractAddress,
    function: functionAbi,
    inputs,
    value: ethValue || '0',
  };
};

/**
 * Prepares a transaction for signing
 */
export const prepareTransaction = async (
  transactionData: TransactionData,
): Promise<PreparedTransaction> => {
  try {
    const encodedData = encodeParameters(
      transactionData.function,
      transactionData.inputs,
    );

    // TODO: Estimate gas
    const estimatedGas = await estimateGas(transactionData);

    return {
      to: transactionData.to,
      data: encodedData,
      value: transactionData.value || '0',
      gas: estimatedGas,
      gasPrice: await getGasPrice(),
    };
  } catch (error) {
    throw new Error(
      `Failed to prepare transaction: ${
        error instanceof Error ? error.message : 'Unknown error'
      }`,
    );
  }
};

/**
 * Estimates gas for a transaction
 */
export const estimateGas = async (
  transactionData: TransactionData,
): Promise<string> => {
  // TODO: Implement actual gas estimation
  // This would typically call eth_estimateGas

  // Placeholder values based on function complexity
  const baseGas = 21000; // Base transaction cost
  const functionGas = transactionData.inputs.length * 5000; // Rough estimate per parameter

  return (baseGas + functionGas).toString();
};

/**
 * Gets current gas price
 */
export const getGasPrice = async (): Promise<string> => {
  // TODO: Implement actual gas price fetching
  // This would typically call eth_gasPrice or use a gas oracle

  // Placeholder: 20 gwei
  return (20 * 1e9).toString();
};

/**
 * Validates transaction inputs
 */
export const validateTransactionInputs = (
  functionAbi: types.Function,
  inputs: TransactionInput[],
): { isValid: boolean; errors: string[] } => {
  const errors: string[] = [];

  // Check all required parameters are provided
  functionAbi.inputs.forEach((param) => {
    const input = inputs.find((i) => i.name === param.name);
    if (!input || !input.value.trim()) {
      if (!param.name.startsWith('_')) {
        // Convention: optional params start with _
        errors.push(`Parameter '${param.name}' is required`);
      }
    }
  });

  // Check for extra parameters
  inputs.forEach((input) => {
    const param = functionAbi.inputs.find((p) => p.name === input.name);
    if (!param) {
      errors.push(`Unknown parameter '${input.name}'`);
    }
  });

  return {
    isValid: errors.length === 0,
    errors,
  };
};
