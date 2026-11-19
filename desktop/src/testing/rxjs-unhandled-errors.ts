import { config } from "rxjs";

/**
 * Runs `fn`, which makes a subscription fail with no error handler (as the app does
 * on purpose wherever the HTTP interceptor reports the failure), and returns the
 * errors RxJS reported as unhandled.
 *
 * RxJS rethrows an unhandled error from its own `setTimeout`, so left alone it lands
 * in whichever later test first waits on a timer, and it reads
 * `config.onUnhandledError` only when that timer fires. Restoring the handler
 * before then, as a plain try/finally does, still leaks. So `fn` runs under Jest's
 * fake timers and the pending timers are flushed while the collecting handler is
 * installed.
 *
 * Keep `fn` synchronous: run any `await` after this returns, on real timers.
 */
export function swallowUnhandledRxjsErrors(fn: () => void): unknown[] {
  const errors: unknown[] = [];
  const previousHandler = config.onUnhandledError;

  jest.useFakeTimers();
  config.onUnhandledError = (error) => errors.push(error);
  try {
    fn();
    jest.runOnlyPendingTimers();
  } finally {
    config.onUnhandledError = previousHandler;
    jest.useRealTimers();
  }

  return errors;
}
