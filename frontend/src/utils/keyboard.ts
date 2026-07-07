const EDITABLE_TAGS = ['INPUT', 'TEXTAREA', 'SELECT', 'BUTTON'];

export function isEditableElement(target: EventTarget | null): boolean {
  if (!(target instanceof HTMLElement)) return false;
  return EDITABLE_TAGS.includes(target.tagName) || target.isContentEditable;
}
