import { TestBed } from '@angular/core/testing';
import { OptionDisplayPipe } from './option-display.pipe';

describe('OptionDisplayPipe', () => {
  let pipe: OptionDisplayPipe;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [OptionDisplayPipe],
      imports: [],
      providers: [OptionDisplayPipe],
    }).compileComponents();

    pipe = TestBed.inject(OptionDisplayPipe);
  });
  it('create an instance', () => {
    expect(pipe).toBeTruthy();
  });

  it('should just return the option', () => {
    const result = pipe.transform('option', []);
    expect(result).toEqual('option');
  });

  it('should just return the option at display key', () => {
    const result = pipe.transform({ test: 'the result' }, [], 'test');
    expect(result).toEqual('the result');
  });

  it('should return the option when there is a value key and display key', () => {
    const options = [
      {
        display: 'best option display',
        value: 'best option value',
      },
    ];
    const result = pipe.transform(options[0], options, 'display', 'value');
    expect(result).toEqual('best option display');
  });
});
