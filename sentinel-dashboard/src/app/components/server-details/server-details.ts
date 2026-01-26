import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { Metric, SentinelService, ServiceStatus } from '../../services/sentinel.service';
import { ChartConfiguration, ChartOptions } from 'chart.js';
import { CommonModule } from '@angular/common';
import { BaseChartDirective, provideCharts, withDefaultRegisterables } from 'ng2-charts';

@Component({
  selector: 'app-server-details',
  imports: [RouterLink,CommonModule,BaseChartDirective],
  providers: [provideCharts(withDefaultRegisterables())],
  templateUrl: './server-details.html',
  styleUrl: './server-details.css',
})
export class ServerDetails implements OnInit {
  serverId: string = '';
  services: ServiceStatus[] = [];

  public lineChartData: ChartConfiguration<'line'>['data'] = {
    labels: [],
    datasets: [
      {
        data: [],
        label: 'CPU Usage (%)',
        fill: true,
        tension: 0.4,
        borderColor: '#3b82f6',
        backgroundColor: 'rgba(59, 130, 246, 0.2)',
        yAxisID: 'y',
      },
      {
        data: [],
        label: 'Memory Usage (MB)',
        fill: true,
        tension: 0.4,
        borderColor: '#8b5cf6',
        backgroundColor: 'rgba(139, 92, 246, 0.2)',
        yAxisID: 'y1',
      }
    ]
  };
  public lineChartOptions: ChartOptions<'line'> = {
    responsive: true,
    maintainAspectRatio: false,
    elements: {
      point: { radius: 0 }
    },
    scales: {
      x: { display: false },
      y: {
        type: 'linear',
        display: true,
        position: 'left',
        title: { display: true, text: 'CPU (%)', color: '#3b82f6' },
        grid: { color: '#334155' },
        ticks: { color: '#94a3b8' }
      },
      y1: {
        type: 'linear',
        display: true,
        position: 'right',
        title: { display: true, text: 'Memory (MB)', color: '#8b5cf6' },
        grid: { drawOnChartArea: false },
        ticks: { color: '#94a3b8' }
      }
    },
    plugins: {
      legend: { labels: { color: '#e2e8f0' } }
    }
  };

  constructor(
    private route: ActivatedRoute,
    private sentinel: SentinelService,
    private cdr: ChangeDetectorRef
  ) { }

  ngOnInit() {
    this.serverId = this.route.snapshot.paramMap.get('id') || '';
    this.refreshData();
    setInterval(() => this.refreshData(), 5000);
  }
  refreshData() {
    this.sentinel.getMetrics(this.serverId).subscribe(metrics => {
      this.updateChart(metrics);
      console.log("Metrics: ", metrics);
      this.cdr.detectChanges();
    });

    this.sentinel.getServiceStatus(this.serverId).subscribe(services => {
      this.services = services || [];
      this.cdr.detectChanges();
    });
  }

  updateChart(metrics: Metric[]) {
    const cpuMetrics = metrics.filter(m => m.metric_type === 'cpu_usage' && m.resource === '').reverse();
    const memMetrics = metrics.filter(m => m.metric_type === 'memory_total_mb' && m.resource === '').reverse();

    this.lineChartData.labels = cpuMetrics.map(m => new Date(m.time).toLocaleTimeString());
    this.lineChartData.datasets[0].data = cpuMetrics.map(m => m.value);
    this.lineChartData.datasets[1].data = memMetrics.map(m => m.value);

    this.lineChartData = { ...this.lineChartData };
  }

}
