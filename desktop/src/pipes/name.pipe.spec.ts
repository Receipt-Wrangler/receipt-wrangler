import { NamePipe } from './name.pipe';

describe('NamePipe', () => {
  it('create an instance', () => {
    const pipe = new NamePipe();
    expect(pipe).toBeTruthy();
  });
});
