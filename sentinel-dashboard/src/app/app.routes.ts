import { Routes } from '@angular/router';
import { Dashboard } from './components/dashboard/dashboard';
import { ServerDetails } from './components/server-details/server-details';

export const routes: Routes = [
     { path: '', component: Dashboard },
     { path: 'server/:id', component: ServerDetails },
];
