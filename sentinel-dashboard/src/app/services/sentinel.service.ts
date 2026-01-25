import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

export interface ServerStatus {
  server_id: string;
  last_seen: string;
  ip_address: string;
}

export interface Metric {
  time: string;
  server_id: string;
  metric_type: string;
  resource: string;
  value: number;
  tags: any;
}

export interface ServiceStatus {
  service_name: string;
  status: number;
  last_seen: string;
}

@Injectable({
  providedIn: 'root',
})
export class SentinelService {
  private apiUrl = 'http://localhost:8080';

  constructor(private http: HttpClient) { }

  getServers(): Observable<ServerStatus[]> {
    return this.http.get<ServerStatus[]>(`${this.apiUrl}/servers`);
  }

  getMetrics(serverId: string): Observable<Metric[]> {
    return this.http.get<Metric[]>(`${this.apiUrl}/metrics/${serverId}`);
  }

  getServiceStatus(serverId: string): Observable<ServiceStatus[]> {
    return this.http.get<ServiceStatus[]>(`${this.apiUrl}/servers/${serverId}/services`);
  }
}
