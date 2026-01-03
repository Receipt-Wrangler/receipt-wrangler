import { InputReadonlyPipe } from './input-readonly.pipe';

describe('InputReadonlyPipe', () => {
  it('create an instance', () => {
    const pipe = new InputReadonlyPipe();
    expect(pipe).toBeTruthy();
  });
});
