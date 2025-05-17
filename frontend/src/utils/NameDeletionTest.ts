// Test script for deleted name records
// This file can be used to verify how deleted records are being processed
import { types } from '@models';

import { createEnhancedName, inspectName } from '../utils/NameTypeUtils';

export function testNameDeletion() {
  console.log('Running name deletion test...');

  // Test case 1: Create a normal name with deleted=true
  const testName1 = {
    address: '0x1234567890123456789012345678901234567890',
    name: 'Test Name 1',
    deleted: true,
    source: 'Test',
    tags: 'Custom',
  };

  // Test case 2: Create a name with deleted="true" (string value)
  const testName2 = {
    address: '0x2345678901234567890123456789012345678901',
    name: 'Test Name 2',
    deleted: 'true',
    source: 'Test',
    tags: 'Custom',
  };

  // Test case 3: Create a name with deleted=1 (number value)
  const testName3 = {
    address: '0x3456789012345678901234567890123456789012',
    name: 'Test Name 3',
    deleted: 1,
    source: 'Test',
    tags: 'Custom',
  };

  // Standard creation method
  const standardName1 = types.Name.createFrom(testName1);
  const standardName2 = types.Name.createFrom(testName2);
  const standardName3 = types.Name.createFrom(testName3);

  // Enhanced creation method
  const enhancedName1 = createEnhancedName(testName1);
  const enhancedName2 = createEnhancedName(testName2);
  const enhancedName3 = createEnhancedName(testName3);

  console.log(
    'Standard creation - deleted as boolean:',
    standardName1.deleted,
    typeof standardName1.deleted,
  );
  console.log(
    'Standard creation - deleted as string:',
    standardName2.deleted,
    typeof standardName2.deleted,
  );
  console.log(
    'Standard creation - deleted as number:',
    standardName3.deleted,
    typeof standardName3.deleted,
  );

  console.log(
    'Enhanced creation - deleted as boolean:',
    enhancedName1.deleted,
    typeof enhancedName1.deleted,
  );
  console.log(
    'Enhanced creation - deleted as string:',
    enhancedName2.deleted,
    typeof enhancedName2.deleted,
  );
  console.log(
    'Enhanced creation - deleted as number:',
    enhancedName3.deleted,
    typeof enhancedName3.deleted,
  );

  // Inspect properties
  console.log('Standard name 1 inspection:', inspectName(standardName1));
  console.log('Enhanced name 1 inspection:', inspectName(enhancedName1));

  // Test boolean comparison
  console.log(
    'StandardName1 deleted === true:',
    standardName1.deleted === true,
  );
  console.log(
    'EnhancedName1 deleted === true:',
    enhancedName1.deleted === true,
  );

  return {
    standardCreation: {
      name1: standardName1,
      name2: standardName2,
      name3: standardName3,
    },
    enhancedCreation: {
      name1: enhancedName1,
      name2: enhancedName2,
      name3: enhancedName3,
    },
  };
}
