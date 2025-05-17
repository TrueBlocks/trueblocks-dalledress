// NameTypeUtils.ts - Helper functions for working with Name types
import { types } from '@models';

/**
 * Enhanced Name creation that ensures the deleted property is a proper boolean
 * @param nameData Original name data from API
 * @returns A properly formatted Name instance with correct deleted property
 */
export function createEnhancedName(
  nameData: types.Name | Record<string, unknown>,
): types.Name {
  // Create a name instance using the standard method
  const name = types.Name.createFrom(nameData);

  // Ensure deleted is a proper boolean value
  if ('deleted' in nameData && nameData.deleted !== undefined) {
    // First convert to boolean
    const isDeleted = Boolean(nameData.deleted);

    // Set the property directly
    name.deleted = isDeleted;

    // Also ensure it's properly defined as a property (for React rendering)
    Object.defineProperty(name, 'deleted', {
      value: isDeleted,
      enumerable: true,
      configurable: true,
      writable: true,
    });
  }

  return name;
}

/**
 * Inspects a name object to verify it has the expected properties
 * Useful for debugging property issues
 */
export function inspectName(name: types.Name): Record<string, unknown> {
  // Get property descriptors
  const descriptors = Object.getOwnPropertyDescriptors(name);

  // Extract key information
  return {
    address: name.address,
    name: name.name,
    deleted: name.deleted,
    deletedType: typeof name.deleted,
    deletedDescriptor: descriptors.deleted || 'not defined',
    ownProps: Object.getOwnPropertyNames(name),
    protoProps: Object.getOwnPropertyNames(Object.getPrototypeOf(name)),
  };
}
