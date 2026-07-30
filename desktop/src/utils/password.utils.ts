/**
 * Secure password generation for the admin-facing "Generate Password" control.
 *
 * Randomness comes from `crypto.getRandomValues` with rejection sampling — a
 * plain `% charsetLength` would skew the distribution toward the first
 * characters of each pool. Generated passwords are shown on screen and handed
 * off by an administrator, so the pools omit visually ambiguous glyphs
 * (`0/O`, `1/l/I`) to keep them readable and transcribable.
 */

const LOWERCASE_CHARACTERS = "abcdefghijkmnopqrstuvwxyz";

const UPPERCASE_CHARACTERS = "ABCDEFGHJKLMNPQRSTUVWXYZ";

const DIGIT_CHARACTERS = "23456789";

const SYMBOL_CHARACTERS = "!@#$%^&*()-_=+[]{}";

/**
 * The character classes a generated password is guaranteed to draw from at
 * least once.
 */
export const PASSWORD_CHARACTER_CLASSES: readonly string[] = [
  LOWERCASE_CHARACTERS,
  UPPERCASE_CHARACTERS,
  DIGIT_CHARACTERS,
  SYMBOL_CHARACTERS,
];

/** Every character a generated password may contain. */
export const PASSWORD_CHARACTER_POOL = PASSWORD_CHARACTER_CLASSES.join("");

/**
 * Default generated password length. Comfortably strong (~120 bits of entropy
 * over the pool above) while staying well inside bcrypt's 72-byte input limit.
 */
export const GENERATED_PASSWORD_LENGTH = 20;

/**
 * Generates a cryptographically random password containing at least one
 * character from each class in {@link PASSWORD_CHARACTER_CLASSES}.
 *
 * @throws RangeError when `length` is too short to satisfy that guarantee.
 */
export function generateSecurePassword(
  length: number = GENERATED_PASSWORD_LENGTH
): string {
  // Guard the loop below: Infinity would spin forever, and NaN would fall
  // straight through and silently yield a password of only the guaranteed
  // characters.
  if (!Number.isInteger(length) || length < PASSWORD_CHARACTER_CLASSES.length) {
    throw new RangeError(
      `Password length must be a whole number of at least ${PASSWORD_CHARACTER_CLASSES.length}.`
    );
  }

  const characters = PASSWORD_CHARACTER_CLASSES.map((characterClass) =>
    randomCharacter(characterClass)
  );

  while (characters.length < length) {
    characters.push(randomCharacter(PASSWORD_CHARACTER_POOL));
  }

  // The class-guaranteed characters were appended in a fixed order, so shuffle
  // to keep their positions unpredictable.
  return shuffle(characters).join("");
}

function randomCharacter(characters: string): string {
  return characters.charAt(randomIndex(characters.length));
}

/**
 * Returns a uniformly distributed integer in `[0, maxExclusive)` by discarding
 * random bytes that fall into the biased tail of the byte range.
 */
function randomIndex(maxExclusive: number): number {
  const buffer = new Uint8Array(1);
  const largestUnbiasedValue = 256 - (256 % maxExclusive);

  for (;;) {
    crypto.getRandomValues(buffer);
    if (buffer[0] < largestUnbiasedValue) {
      return buffer[0] % maxExclusive;
    }
  }
}

function shuffle(characters: string[]): string[] {
  for (let i = characters.length - 1; i > 0; i--) {
    const j = randomIndex(i + 1);
    [characters[i], characters[j]] = [characters[j], characters[i]];
  }

  return characters;
}
