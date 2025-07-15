import { useWalletConnectContext } from '@contexts';
import { useWallet } from '@hooks';

import { PreparedTransaction } from './transactionBuilder';

export interface WalletConnectionProps {
  onTransactionSigned?: (txHash: string) => void;
  onError?: (error: string) => void;
}

export const useWalletConnection = ({
  onTransactionSigned,
  onError,
}: WalletConnectionProps = {}) => {
  const { session } = useWalletConnectContext();
  const { walletAddress, walletChainId, isWalletConnected } = useWallet();

  /**
   * Send a transaction via WalletConnect for signing
   * This is a placeholder implementation - the actual WalletConnect transaction
   * sending would depend on the specific WalletConnect SDK version and setup
   */
  const sendTransaction = async (
    preparedTx: PreparedTransaction,
  ): Promise<string> => {
    if (!isWalletConnected || !walletAddress) {
      throw new Error('Wallet not connected');
    }

    if (!session.isConnected) {
      throw new Error('WalletConnect session not available');
    }

    try {
      // Prepare the transaction request for WalletConnect
      const transactionRequest = {
        from: walletAddress,
        to: preparedTx.to,
        data: preparedTx.data,
        value: preparedTx.value,
        gas: preparedTx.gas,
        gasPrice: preparedTx.gasPrice,
      };

      // TODO: Implement actual WalletConnect transaction sending
      // This would use the WalletConnect client to send the transaction
      // to the connected wallet for signing
      console.log('Sending transaction via WalletConnect:', transactionRequest);

      // For now, simulate a transaction hash
      // In reality, this would come from the WalletConnect response
      const txHash = `0x${Math.random().toString(16).substring(2, 66)}`;

      if (onTransactionSigned) {
        onTransactionSigned(txHash);
      }

      return txHash;
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : 'Transaction failed';

      if (onError) {
        onError(errorMessage);
      }

      throw new Error(errorMessage);
    }
  };

  /**
   * Sign a message via WalletConnect
   */
  const signMessage = async (message: string): Promise<string> => {
    if (!isWalletConnected || !walletAddress) {
      throw new Error('Wallet not connected');
    }

    if (!session.isConnected) {
      throw new Error('WalletConnect session not available');
    }

    try {
      // TODO: Implement actual WalletConnect message signing
      console.log('Signing message via WalletConnect:', message);

      // For now, simulate a signature
      const signature = `0x${Math.random().toString(16).substring(2, 130)}`;

      return signature;
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : 'Message signing failed';

      if (onError) {
        onError(errorMessage);
      }

      throw new Error(errorMessage);
    }
  };

  /**
   * Get the current wallet info
   */
  const getWalletInfo = () => ({
    address: walletAddress,
    chainId: walletChainId,
    isConnected: isWalletConnected,
    session: session,
  });

  return {
    // Transaction methods
    sendTransaction,
    signMessage,

    // Wallet info
    getWalletInfo,

    // Connection status
    isConnected: isWalletConnected,
    address: walletAddress,
    chainId: walletChainId,
  };
};

// Export the component for backwards compatibility
export const WalletConnection: React.FC<WalletConnectionProps> = (props) => {
  const _wallet = useWalletConnection(props);

  // This component doesn't render anything - it's just for the hook
  return null;
};
