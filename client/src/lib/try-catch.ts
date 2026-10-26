// import type { ApiErrorResponse } from "./requests";

// type Success<T> = readonly [T, null];
// type Failure<E> = readonly [null, E];
// type RequestFailure = readonly [null, ApiErrorResponse];
// type Result<T, E = Error> = Success<T> | Failure<E>;
// export type RequestResult<T> = Success<T> | RequestFailure;

// export async function tryCatch<T, E = Error>(
//   promise: Promise<T>,
//   finallyCallback?: () => void,
// ): Promise<Result<T, E>> {
//   try {
//     const data = await promise;
//     return [data, null];
//   } catch (error) {
//     return [null, error as E] as const;
//   } finally {
//     finallyCallback?.();
//   }
// }

// export async function requestTryCatch<T>(
//   promise: Promise<T>,
//   finallyCallback?: () => void,
// ): Promise<RequestResult<T>> {
//   try {
//     const data = await promise;
//     return [data, null];
//   } catch (error) {
//     return [null, error as ApiErrorResponse] as const;
//   } finally {
//     finallyCallback?.();
//   }
// }

// export async function returnTryCatchError<T>(
//   data: ApiErrorResponse,
//   finallyCallback?: () => void,
// ): Promise<RequestResult<T>> {
//   try {
//     return Promise.resolve([null, data] as const);
//   } catch {
//     return Promise.resolve([
//       null,
//       {
//         statusCode: 500,
//         message: "An unexpected error occurred",
//       },
//     ] as const);
//   } finally {
//     finallyCallback?.();
//   }
// }
