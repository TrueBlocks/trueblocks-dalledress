import { FieldRenderer, FormField } from '@components';
import { types } from '@models';
import { getDisplayAddress } from '@utils';

function displayHash(v: unknown): string {
  try {
    if (!v) return '';
    if (typeof v === 'string') {
      const s = v.startsWith('0x') ? v : `0x${v}`;
      if (s.length <= 18) return s;
      return `${s.slice(0, 10)}…${s.slice(-8)}`;
    }
    const obj = v as { hash?: number[] };
    const arr = obj && Array.isArray(obj.hash) ? obj.hash : undefined;
    if (arr && arr.length > 0) {
      const hex = `0x${arr.map((b) => b.toString(16).padStart(2, '0')).join('')}`;
      if (hex.length <= 18) return hex;
      return `${hex.slice(0, 10)}…${hex.slice(-8)}`;
    }
    const s = String(v);
    if (s.length <= 18) return s;
    return `${s.slice(0, 10)}…${s.slice(-8)}`;
  } catch {
    return String(v || '');
  }
}

/**
 * Custom bank statement-style detail panel for Statement rows
 * Displays transaction data in a professional bank statement format with ether values
 */
export const renderStatementDetailPanel = (
  rowData: Record<string, unknown> | null,
) => {
  if (!rowData) return null;
  const statement = rowData as unknown as types.Statement;

  // Helper function to format values as ether
  const formatValue = (value: unknown) => {
    if (!value || value === '0') return '-';

    // Convert wei to ether using the statement's decimals
    const valueStr = typeof value === 'string' ? value : String(value);
    const decimals = statement.decimals || 18; // Default to 18 for ETH

    try {
      // Simple conversion: divide by 10^decimals
      const divisor = Math.pow(10, decimals);
      const etherValue = parseFloat(valueStr) / divisor;

      // Format to 3 decimal places
      return etherValue.toFixed(3);
    } catch {
      return valueStr; // Fallback to original value if conversion fails
    }
  };

  // Check if reconciled (no correcting reasons means it's reconciled)
  const isReconciled =
    !statement.correctingReasons || statement.correctingReasons === '';

  return (
    <div
      style={{
        fontFamily: 'monospace',
        fontSize: '14px',
        lineHeight: '1.6',
      }}
    >
      <div
        style={{
          textAlign: 'center',
          marginBottom: '20px',
          borderBottom: '2px solid #333',
          paddingBottom: '10px',
        }}
      >
        <h3>Jay&apos;s Bank Statement</h3>
        <div style={{ fontSize: '12px', color: '#666' }}>
          Block: {statement.blockNumber} | Tx: {statement.transactionIndex} |
          Date: {new Date(statement.timestamp * 1000).toLocaleDateString()}
        </div>
      </div>

      <div style={{ marginBottom: '15px' }}>
        <strong>Transaction Identification:</strong>
        <br />
        <div style={{ marginLeft: '10px' }}>
          Account: {getDisplayAddress(statement.accountedFor)}
          <br />
          Asset: {statement.symbol} ({getDisplayAddress(statement.asset)})
          <br />
          From: {getDisplayAddress(statement.sender)}
          <br />
          To: {getDisplayAddress(statement.recipient)}
        </div>
      </div>

      <table
        style={{
          width: '100%',
          borderCollapse: 'collapse',
          marginBottom: '15px',
        }}
      >
        <thead>
          <tr style={{ borderBottom: '1px solid #333' }}>
            <th style={{ textAlign: 'left', padding: '5px' }}>Type</th>
            <th style={{ textAlign: 'right', padding: '5px' }}>In</th>
            <th style={{ textAlign: 'right', padding: '5px' }}>Out</th>
            <th style={{ textAlign: 'right', padding: '5px' }}>Gas</th>
            <th style={{ textAlign: 'right', padding: '5px' }}>Balance</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td style={{ padding: '5px' }}>Beginning Balance</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>
              {formatValue(statement.begBal)}
            </td>
          </tr>
          <tr>
            <td style={{ padding: '5px' }}>Regular Amount</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>
              {formatValue(statement.amountIn)}
            </td>
            <td style={{ textAlign: 'right', padding: '5px' }}>
              {formatValue(statement.amountOut)}
            </td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
          </tr>
          <tr>
            <td style={{ padding: '5px' }}>Internal Transfers</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>
              {formatValue(statement.internalIn)}
            </td>
            <td style={{ textAlign: 'right', padding: '5px' }}>
              {formatValue(statement.internalOut)}
            </td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
          </tr>
          <tr>
            <td style={{ padding: '5px' }}>Self Destruct</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>
              {formatValue(statement.selfDestructIn)}
            </td>
            <td style={{ textAlign: 'right', padding: '5px' }}>
              {formatValue(statement.selfDestructOut)}
            </td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
          </tr>
          <tr>
            <td style={{ padding: '5px' }}>Gas Costs</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>
              {formatValue(statement.gasOut)}
            </td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
          </tr>
          <tr>
            <td style={{ padding: '5px' }}>Miner Rewards</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>
              {formatValue(statement.minerBaseRewardIn)}
            </td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
          </tr>
          <tr style={{ borderTop: '1px solid #333', fontWeight: 'bold' }}>
            <td style={{ padding: '5px' }}>Ending Balance</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>-</td>
            <td style={{ textAlign: 'right', padding: '5px' }}>
              {formatValue(statement.endBal)}
            </td>
          </tr>
        </tbody>
      </table>

      <div
        style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          marginTop: '10px',
        }}
      >
        <span style={{ marginRight: '10px' }}>Reconciled:</span>
        {isReconciled ? (
          <span style={{ color: 'green', fontSize: '18px' }}>✓</span>
        ) : (
          <span style={{ color: 'red', fontSize: '18px' }}>✗</span>
        )}
        {!isReconciled && statement.correctingReasons && (
          <span style={{ marginLeft: '10px', fontSize: '12px', color: '#666' }}>
            ({statement.correctingReasons})
          </span>
        )}
      </div>
    </div>
  );
};

export const renderReceiptDetailPanel = (
  rowData: Record<string, unknown> | null,
) => {
  if (!rowData) return null;
  const receipt = rowData as unknown as types.Receipt;

  const has = (v: unknown) => v !== undefined && v !== null && v !== '';

  let gasCostWei: string | undefined;
  try {
    if (has(receipt.gasUsed) && has(receipt.effectiveGasPrice)) {
      const gu = BigInt(String(receipt.gasUsed));
      const egp = BigInt(String(receipt.effectiveGasPrice));
      gasCostWei = (gu * egp).toString();
    }
  } catch {
    gasCostWei = undefined;
  }

  const sections: {
    title: string;
    fields: FormField<Record<string, unknown>>[];
  }[] = [
    {
      title: 'Status',
      fields: [
        has(receipt.status)
          ? {
              key: 'status',
              name: 'status',
              label: 'Status',
              type: 'value',
              value: String(receipt.status),
              readOnly: true,
            }
          : undefined,
        has(receipt.isError)
          ? {
              key: 'isError',
              name: 'isError',
              label: 'Error',
              type: 'value',
              value: String(Boolean(receipt.isError)),
              readOnly: true,
            }
          : undefined,
        Array.isArray(receipt.logs)
          ? {
              key: 'logCount',
              name: 'logCount',
              label: 'Log Count',
              type: 'value',
              value: String(receipt.logs.length),
              readOnly: true,
            }
          : undefined,
      ].filter(Boolean) as FormField<Record<string, unknown>>[],
    },
    {
      title: 'Gas',
      fields: [
        has(receipt.gasUsed)
          ? {
              key: 'gasUsed',
              name: 'gasUsed',
              label: 'Gas Used',
              type: 'gas',
              value: String(receipt.gasUsed),
              readOnly: true,
            }
          : undefined,
        has(receipt.cumulativeGasUsed)
          ? {
              key: 'cumulativeGasUsed',
              name: 'cumulativeGasUsed',
              label: 'Cumulative Gas Used',
              type: 'gas',
              value: String(receipt.cumulativeGasUsed),
              readOnly: true,
            }
          : undefined,
        has(receipt.effectiveGasPrice)
          ? {
              key: 'effectiveGasPrice',
              name: 'effectiveGasPrice',
              label: 'Effective Gas Price',
              type: 'gas',
              value: String(receipt.effectiveGasPrice),
              readOnly: true,
            }
          : undefined,
        gasCostWei
          ? {
              key: 'gasCost',
              name: 'gasCost',
              label: 'Gas Cost',
              type: 'ether',
              value: gasCostWei,
              readOnly: true,
            }
          : undefined,
      ].filter(Boolean) as FormField<Record<string, unknown>>[],
    },
    {
      title: 'Participants',
      fields: [
        has(receipt.from)
          ? {
              key: 'from',
              name: 'from',
              label: 'From',
              type: 'address',
              value: getDisplayAddress(receipt.from),
              readOnly: true,
            }
          : undefined,
        has(receipt.to)
          ? {
              key: 'to',
              name: 'to',
              label: 'To',
              type: 'address',
              value: getDisplayAddress(receipt.to),
              readOnly: true,
            }
          : undefined,
        has(receipt.contractAddress)
          ? {
              key: 'contractAddress',
              name: 'contractAddress',
              label: 'Contract',
              type: 'address',
              value: getDisplayAddress(receipt.contractAddress),
              readOnly: true,
            }
          : undefined,
      ].filter(Boolean) as FormField<Record<string, unknown>>[],
    },

    {
      title: 'References',
      fields: [
        has(receipt.blockNumber)
          ? {
              key: 'blockNumber',
              name: 'blockNumber',
              label: 'Block Number',
              type: 'blknum',
              value: String(receipt.blockNumber),
              readOnly: true,
            }
          : undefined,
        has(receipt.transactionIndex)
          ? {
              key: 'transactionIndex',
              name: 'transactionIndex',
              label: 'Tx Index',
              type: 'txnum',
              value: String(receipt.transactionIndex),
              readOnly: true,
            }
          : undefined,
        has(receipt.transactionHash)
          ? {
              key: 'transactionHash',
              name: 'transactionHash',
              label: 'Tx Hash',
              type: 'hash',
              value: displayHash(receipt.transactionHash),
              readOnly: true,
            }
          : undefined,
        has(receipt.blockHash)
          ? {
              key: 'blockHash',
              name: 'blockHash',
              label: 'Block Hash',
              type: 'hash',
              value: displayHash(receipt.blockHash),
              readOnly: true,
            }
          : undefined,
      ].filter(Boolean) as FormField<Record<string, unknown>>[],
    },
  ];

  return (
    <div
      style={{
        fontFamily: 'monospace',
        fontSize: '14px',
        lineHeight: '1.6',
      }}
    >
      <div
        style={{
          textAlign: 'center',
          marginBottom: '10px',
          borderBottom: '2px solid #333',
          paddingBottom: '10px',
        }}
      >
        <h3>Transaction Receipt</h3>
        <div style={{ fontSize: '12px', color: '#666' }}>
          {`Block: ${String(receipt.blockNumber)} | Tx: ${String(
            receipt.transactionIndex,
          )}`}
        </div>
      </div>

      <table
        style={{
          width: '100%',
          borderCollapse: 'collapse',
        }}
      >
        <tbody>
          {sections.map((section) => (
            <>
              <tr key={`${section.title}-header`}>
                <td
                  colSpan={2}
                  style={{
                    padding: '8px 6px',
                    background: 'var(--mantine-color-gray-0)',
                    fontWeight: 600,
                  }}
                >
                  {section.title}
                </td>
              </tr>
              {section.fields.map((field, idx) => (
                <tr key={`${section.title}-${idx}`}>
                  <td
                    style={{
                      width: '220px',
                      padding: '4px 6px',
                      verticalAlign: 'top',
                      color: 'var(--mantine-color-dimmed)',
                      whiteSpace: 'nowrap',
                    }}
                  >
                    {field.label}
                  </td>
                  <td style={{ padding: '4px 6px' }}>
                    <FieldRenderer field={field} mode="display" tableCell />
                  </td>
                </tr>
              ))}
              <tr key={`${section.title}-separator`}>
                <td
                  colSpan={2}
                  style={{ borderBottom: '1px solid #333', padding: 0 }}
                />
              </tr>
            </>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export const renderTraceDetailPanel = (
  rowData: Record<string, unknown> | null,
) => {
  if (!rowData) return null;
  const trace = rowData as unknown as types.Trace;

  const has = (v: unknown) => v !== undefined && v !== null && v !== '';
  const joinAddr = (arr?: number[]) =>
    Array.isArray(arr) && arr.length > 0 ? arr.join('.') : '0';

  const sections: {
    title: string;
    fields: FormField<Record<string, unknown>>[];
  }[] = [
    {
      title: 'Status',
      fields: [
        has(trace.type)
          ? {
              key: 'type',
              name: 'type',
              label: 'Type',
              type: 'value',
              value: String(trace.type),
              readOnly: true,
            }
          : undefined,
        has(trace.error)
          ? {
              key: 'error',
              name: 'error',
              label: 'Error',
              type: 'value',
              value: String(trace.error),
              readOnly: true,
            }
          : undefined,
        has(trace.subtraces)
          ? {
              key: 'subtraces',
              name: 'subtraces',
              label: 'Subtraces',
              type: 'value',
              value: String(trace.subtraces),
              readOnly: true,
            }
          : undefined,
      ].filter(Boolean) as FormField<Record<string, unknown>>[],
    },
    {
      title: 'Gas',
      fields: [
        has(trace.action?.gas)
          ? {
              key: 'actionGas',
              name: 'actionGas',
              label: 'Action Gas',
              type: 'gas',
              value: String(trace.action?.gas),
              readOnly: true,
            }
          : undefined,
        has(trace.result?.gasUsed)
          ? {
              key: 'resultGasUsed',
              name: 'resultGasUsed',
              label: 'Result Gas Used',
              type: 'gas',
              value: String(trace.result?.gasUsed),
              readOnly: true,
            }
          : undefined,
      ].filter(Boolean) as FormField<Record<string, unknown>>[],
    },
    {
      title: 'Participants',
      fields: [
        has(trace.action?.from)
          ? {
              key: 'from',
              name: 'from',
              label: 'From',
              type: 'address',
              value: getDisplayAddress(trace.action?.from),
              readOnly: true,
            }
          : undefined,
        has(trace.action?.to)
          ? {
              key: 'to',
              name: 'to',
              label: 'To',
              type: 'address',
              value: getDisplayAddress(trace.action?.to),
              readOnly: true,
            }
          : undefined,
        has(trace.action?.address)
          ? {
              key: 'address',
              name: 'address',
              label: 'Contract',
              type: 'address',
              value: getDisplayAddress(trace.action?.address),
              readOnly: true,
            }
          : undefined,
        has(trace.action?.author)
          ? {
              key: 'author',
              name: 'author',
              label: 'Author',
              type: 'address',
              value: getDisplayAddress(trace.action?.author),
              readOnly: true,
            }
          : undefined,
        has(trace.action?.refundAddress)
          ? {
              key: 'refundAddress',
              name: 'refundAddress',
              label: 'Refund',
              type: 'address',
              value: getDisplayAddress(trace.action?.refundAddress),
              readOnly: true,
            }
          : undefined,
        has(trace.action?.selfDestructed)
          ? {
              key: 'selfDestructed',
              name: 'selfDestructed',
              label: 'Self Destructed',
              type: 'address',
              value: getDisplayAddress(trace.action?.selfDestructed),
              readOnly: true,
            }
          : undefined,
      ].filter(Boolean) as FormField<Record<string, unknown>>[],
    },
    {
      title: 'References',
      fields: [
        has(trace.blockNumber)
          ? {
              key: 'blockNumber',
              name: 'blockNumber',
              label: 'Block Number',
              type: 'blknum',
              value: String(trace.blockNumber),
              readOnly: true,
            }
          : undefined,
        has(trace.transactionIndex)
          ? {
              key: 'transactionIndex',
              name: 'transactionIndex',
              label: 'Tx Index',
              type: 'txnum',
              value: String(trace.transactionIndex),
              readOnly: true,
            }
          : undefined,
        has(trace.traceAddress)
          ? {
              key: 'traceAddress',
              name: 'traceAddress',
              label: 'Trace Address',
              type: 'value',
              value: joinAddr(trace.traceAddress),
              readOnly: true,
            }
          : undefined,
        has(trace.transactionHash)
          ? {
              key: 'transactionHash',
              name: 'transactionHash',
              label: 'Tx Hash',
              type: 'hash',
              value: displayHash(trace.transactionHash),
              readOnly: true,
            }
          : undefined,
        has(trace.blockHash)
          ? {
              key: 'blockHash',
              name: 'blockHash',
              label: 'Block Hash',
              type: 'hash',
              value: displayHash(trace.blockHash),
              readOnly: true,
            }
          : undefined,
      ].filter(Boolean) as FormField<Record<string, unknown>>[],
    },
  ];

  return (
    <div
      style={{
        fontFamily: 'monospace',
        fontSize: '14px',
        lineHeight: '1.6',
      }}
    >
      <div
        style={{
          textAlign: 'center',
          marginBottom: '10px',
          borderBottom: '2px solid #333',
          paddingBottom: '10px',
        }}
      >
        <h3>Transaction Trace</h3>
        <div style={{ fontSize: '12px', color: '#666' }}>
          {`Block: ${String(trace.blockNumber)} | Tx: ${String(
            trace.transactionIndex,
          )}`}
        </div>
      </div>

      <table
        style={{
          width: '100%',
          borderCollapse: 'collapse',
        }}
      >
        <tbody>
          {sections.map((section) => (
            <>
              <tr key={`${section.title}-header`}>
                <td
                  colSpan={2}
                  style={{
                    padding: '8px 6px',
                    background: 'var(--mantine-color-gray-0)',
                    fontWeight: 600,
                  }}
                >
                  {section.title}
                </td>
              </tr>
              {section.fields.map((field, idx) => (
                <tr key={`${section.title}-${idx}`}>
                  <td
                    style={{
                      width: '220px',
                      padding: '4px 6px',
                      verticalAlign: 'top',
                      color: 'var(--mantine-color-dimmed)',
                      whiteSpace: 'nowrap',
                    }}
                  >
                    {field.label}
                  </td>
                  <td style={{ padding: '4px 6px' }}>
                    <FieldRenderer field={field} mode="display" tableCell />
                  </td>
                </tr>
              ))}
              <tr key={`${section.title}-separator`}>
                <td
                  colSpan={2}
                  style={{ borderBottom: '1px solid #333', padding: 0 }}
                />
              </tr>
            </>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export const renderLogDetailPanel = (
  rowData: Record<string, unknown> | null,
) => {
  if (!rowData) return null;
  const log = rowData as unknown as types.Log;

  const has = (v: unknown) => v !== undefined && v !== null && v !== '';
  const decoded = Boolean(log.articulatedLog);
  const topicsStr = Array.isArray(log.topics)
    ? log.topics.map((t) => displayHash(t)).join(', ')
    : '';

  const sections: {
    title: string;
    fields: FormField<Record<string, unknown>>[];
  }[] = [
    {
      title: 'Log Entry',
      fields: [
        has(log.address)
          ? {
              key: 'address',
              name: 'address',
              label: 'Address',
              type: 'address',
              value: getDisplayAddress(log.address),
              readOnly: true,
            }
          : undefined,
        has(topicsStr)
          ? {
              key: 'topics',
              name: 'topics',
              label: 'Topics',
              type: 'text',
              value: topicsStr,
              readOnly: true,
            }
          : undefined,
        has(log.data)
          ? {
              key: 'data',
              name: 'data',
              label: 'Data',
              type: 'bytes',
              value: displayHash(log.data as unknown as string),
              readOnly: true,
            }
          : undefined,
        {
          key: 'articulatedLog',
          name: 'articulatedLog',
          label: 'Articulated Log',
          type: 'text',
          value: decoded ? 'decoded' : 'N/A',
          readOnly: true,
        },
        {
          key: 'decoded',
          name: 'decoded',
          label: 'Decoded',
          type: 'text',
          customRender: (
            <span style={{ fontSize: '16px' }}>{decoded ? '✅' : '✗'}</span>
          ),
          readOnly: true,
        },
      ].filter(Boolean) as FormField<Record<string, unknown>>[],
    },
  ];

  return (
    <div
      style={{
        fontFamily: 'monospace',
        fontSize: '14px',
        lineHeight: '1.6',
      }}
    >
      <div
        style={{
          textAlign: 'center',
          marginBottom: '10px',
          borderBottom: '2px solid #333',
          paddingBottom: '10px',
        }}
      >
        <h3>Log Entry</h3>
        <div style={{ fontSize: '12px', color: '#666' }}>
          {`Block: ${String(log.blockNumber)} | Tx: ${String(
            log.transactionIndex,
          )} | Log Index: ${String(log.logIndex)}`}
        </div>
      </div>

      <table
        style={{
          width: '100%',
          borderCollapse: 'collapse',
        }}
      >
        <tbody>
          {sections.map((section) => (
            <>
              <tr key={`${section.title}-header`}>
                <td
                  colSpan={2}
                  style={{
                    padding: '8px 6px',
                    background: 'var(--mantine-color-gray-0)',
                    fontWeight: 600,
                  }}
                >
                  {section.title}
                </td>
              </tr>
              {section.fields.map((field, idx) => (
                <tr key={`${section.title}-${idx}`}>
                  <td
                    style={{
                      width: '220px',
                      padding: '4px 6px',
                      verticalAlign: 'top',
                      color: 'var(--mantine-color-dimmed)',
                      whiteSpace: 'nowrap',
                    }}
                  >
                    {field.label}
                  </td>
                  <td style={{ padding: '4px 6px' }}>
                    <FieldRenderer field={field} mode="display" tableCell />
                  </td>
                </tr>
              ))}
              <tr key={`${section.title}-separator`}>
                <td
                  colSpan={2}
                  style={{ borderBottom: '1px solid #333', padding: 0 }}
                />
              </tr>
            </>
          ))}
        </tbody>
      </table>
    </div>
  );
};

/**
 * Configuration object for exports view detail panels
 * Maps data facets to their corresponding detail panel functions
 */
export const exportsDetailPanels = {
  'exports.statements': renderStatementDetailPanel,
  'exports.receipts': renderReceiptDetailPanel,
  'exports.traces': renderTraceDetailPanel,
  'exports.logs': renderLogDetailPanel,
};
