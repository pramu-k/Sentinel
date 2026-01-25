import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { SentinelService, ServerStatus } from '../../services/sentinel.service';
import { RouterLink } from '@angular/router';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-dashboard',
  imports: [CommonModule,RouterLink],
  templateUrl: './dashboard.html',
  styleUrl: './dashboard.css',
})
export class Dashboard implements OnInit {
  servers: ServerStatus[] = [];
    errorMessage: string = '';

    constructor(private sentinel: SentinelService, private cdr: ChangeDetectorRef) { }

    ngOnInit() {
        this.refreshServers();
        setInterval(() => this.refreshServers(), 5000);
    }

    refreshServers() {
        this.sentinel.getServers().subscribe({
            next: (data) => {
                this.servers = data || [];
                this.errorMessage = '';
                console.log("Servers Type:", typeof this.servers, "Is Array:", Array.isArray(this.servers), "Length:", this.servers.length);
                console.log("Servers Data:", this.servers);
                this.cdr.detectChanges();
            },
            error: (err: any) => {
                console.error("API Error:", err);
                this.errorMessage = `Failed to load data: ${err.message || 'Unknown error'}`;
                this.servers = [];
                this.cdr.detectChanges();
            }
        });
    }

    isAlive(lastSeen: string): boolean {
        const last = new Date(lastSeen).getTime();
        const now = new Date().getTime();
        return (now - last) < 10000;
    }

}
