import { Routes } from '@angular/router';
import { LoginComponent } from './auth/login/login.component';
import { RegisterComponent } from './auth/register/register.component';
import { AppComponent } from './app.component';
import { DefaultPageComponent } from './default-page/default-page.component';

export const routes: Routes = [
    { path: '', component: DefaultPageComponent },
    { path: 'login', component: LoginComponent },
    { path: 'register', component: RegisterComponent }
];
