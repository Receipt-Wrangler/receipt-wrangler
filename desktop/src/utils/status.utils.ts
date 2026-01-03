export function formatStatus(status: string): string {
  if (!status) return '';
  let result = status.toLowerCase();
  const parts = result.split('_');
  const words: string[] = [];

  parts.forEach((part) => {
    const letters = part.split('');
    letters[0] = letters[0].toUpperCase();
    words.push(letters.join(''));
  });

  return words.join(' ');
}
