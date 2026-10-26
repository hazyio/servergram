
export class Logger {
  public static info(message: unknown, ...props: unknown[]): void {
    console.info(`[INFO] ${message}`, ...props);
  }

  public static warn(message: string): void {
    console.warn(`[WARN] ${message}`);
  }

  public static error(message: unknown, ...props: unknown[]): void {
    if (import.meta.env.DEV) {
      console.error(`[ERROR] ${message}`, ...props);
    } else {
      // In production, we might want to send this to a logging service
      // or handle it differently, but for now, we just log it.
      //   Sentry.captureMessage(message);
    }
  }

  public static exception(error: Error, ...props: unknown[]): void {
    if (import.meta.env.DEV) {
      console.error(`[EXCEPTION] ${error.message}`, ...props);
    } else {
      //   Sentry.captureException(error);
    }
  }

  public static debug(message: unknown, ...props: unknown[]): void {
    if (import.meta.env.DEV) {
      console.debug(`[DEBUG] ${message}`, ...props);
    }
  }
}
