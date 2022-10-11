import { ClassProvider } from '@angular/core';
import { ApiClient } from '@msx/http';
import { Observable } from 'rxjs';

export class MockApiClient {
	request: <T>() => Observable<T>;

	constructor() {
		this.request = jest.fn();
	}
}

export const MockApiProvider: ClassProvider = {
	provide: ApiClient,
	useClass: MockApiClient
};
