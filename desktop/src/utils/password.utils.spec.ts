import {
  GENERATED_PASSWORD_LENGTH,
  PASSWORD_CHARACTER_CLASSES,
  PASSWORD_CHARACTER_POOL,
  generateSecurePassword,
} from "./password.utils";

const ITERATIONS = 200;

describe("generateSecurePassword", () => {
  it("defaults to the configured length", () => {
    expect(generateSecurePassword().length).toEqual(GENERATED_PASSWORD_LENGTH);
  });

  it("honors an explicit length", () => {
    expect(generateSecurePassword(32).length).toEqual(32);
  });

  it("throws when the length cannot satisfy every character class", () => {
    expect(() =>
      generateSecurePassword(PASSWORD_CHARACTER_CLASSES.length - 1)
    ).toThrow(RangeError);
  });

  // Without the integer guard, Infinity spins the fill loop forever and NaN
  // falls through it, returning only the class-guaranteed characters.
  it.each([Infinity, -Infinity, NaN, 20.5])(
    "throws for the non-integer length %p",
    (length) => {
      expect(() => generateSecurePassword(length)).toThrow(RangeError);
    }
  );

  it("always contains at least one character from each class", () => {
    for (let i = 0; i < ITERATIONS; i++) {
      const password = generateSecurePassword();

      PASSWORD_CHARACTER_CLASSES.forEach((characterClass) => {
        expect(
          [...password].some((character) => characterClass.includes(character))
        ).toEqual(true);
      });
    }
  });

  it("only uses characters from the pool", () => {
    for (let i = 0; i < ITERATIONS; i++) {
      [...generateSecurePassword()].forEach((character) => {
        expect(PASSWORD_CHARACTER_POOL).toContain(character);
      });
    }
  });

  it("excludes visually ambiguous characters", () => {
    [..."0O1lI"].forEach((character) => {
      expect(PASSWORD_CHARACTER_POOL).not.toContain(character);
    });
  });

  it("does not repeat itself", () => {
    const passwords = new Set<string>();

    for (let i = 0; i < ITERATIONS; i++) {
      passwords.add(generateSecurePassword());
    }

    expect(passwords.size).toEqual(ITERATIONS);
  });

  it("sources randomness from crypto and not Math.random", () => {
    const getRandomValues = jest.spyOn(crypto, "getRandomValues");
    const random = jest.spyOn(Math, "random");

    try {
      generateSecurePassword();

      expect(getRandomValues).toHaveBeenCalled();
      expect(random).not.toHaveBeenCalled();
    } finally {
      getRandomValues.mockRestore();
      random.mockRestore();
    }
  });

  it("rejects biased random bytes instead of taking them modulo the pool size", () => {
    // The first draw is from the 25-character lowercase class, whose unbiased
    // range is 0-249 — bytes 250-255 must be discarded and redrawn. A naive
    // `% 25` implementation would consume them instead of rejecting them.
    const scriptedBytes = [250, 251, 252, 253, 254, 255];
    let call = 0;
    const getRandomValues = jest
      .spyOn(crypto, "getRandomValues")
      .mockImplementation((array: any) => {
        array[0] = call < scriptedBytes.length ? scriptedBytes[call] : 5;
        call++;
        return array;
      });

    try {
      const password = generateSecurePassword(4);

      // Every scripted byte was rejected, so the first draw fell through to 5.
      expect(call).toBeGreaterThan(scriptedBytes.length);
      expect(password.length).toEqual(4);
    } finally {
      getRandomValues.mockRestore();
    }
  });
});
