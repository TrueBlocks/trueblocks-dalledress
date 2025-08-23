import { useEffect, useState } from 'react';

import { GetProjectAddress } from '@app';
import { TabView } from '@layout';
import { Text } from '@mantine/core';
import { base } from '@models';
import { addressToHex } from '@utils';

export const DalleDress = () => {
  const [address, setAddress] = useState('');
  useEffect(() => {
    GetProjectAddress().then((addr: base.Address) => {
      if (addr && addr.address && addr.address.length > 0) {
        setAddress(addressToHex(addr));
      }
    });
  }, []);
  const tabs = [
    {
      label: 'Template1',
      value: 'template1',
      content: <Text>AddressEntry: {address}</Text>,
    },
  ];
  return <TabView tabs={tabs} route="dalledress" />;
};
