import { Injectable } from '@angular/core';
import { MockApiClient } from './api-client';
import { Observable, of } from 'rxjs';

@Injectable()
export class MockMonitorService {

  constructor(public apiClient: MockApiClient) { }
  public createDeviceMetricTemplate(_: any, standardValues: any): Observable<any> {
    return of(new Object());
  };
  public removeDeviceMetricTemplate(_: any): void {};
  public provisionNSO(_: string): void {};
  deprovisionNSO(_: string): void {};
  startCollecting(_: string): void {};
  addMetricType(_: string): void {};
  addQueryTemplateForService(_: any): void {};
  addServiceIdToQuery(_: string): void {};
  removeServiceFromQuery(_: string): void {};
  getNotifications(_: string): void {};
}
