export const MSG_RETRY = 'Retry'
export const MSG_ADD_PERSON = 'Add person'
export const MSG_SAVING = 'Saving…'
export const MSG_SAVED = 'Saved'
export const MSG_LOAD_FAILED = 'Failed to load the bill.'
export const MSG_NETWORK = 'Could not reach the server. Check your connection and try again.'
export const MSG_TIMEOUT = 'The request took too long. Try again.'
export const MSG_INTERNAL = 'Something went wrong. Try again.'
export const MSG_NAME_REQUIRED = 'Name is required.'
export const MSG_DUPLICATE_NAME = 'Name must be unique.'
export const MSG_PERCENTAGE_OUT_OF_RANGE = 'Percentage must be greater than 0 and at most 100.00.'
export const MSG_AT_LEAST_ONE_SHARE = 'At least one share is required.'

export function nameTooLongMessage(max: number): string {
  return `Name must be at most ${max} characters.`
}

export function tooManySharesMessage(max: number): string {
  return `At most ${max} shares are allowed.`
}

export function sumMismatchMessage(got: string): string {
  return `Percentages must sum to 100.00 (got ${got}).`
}
