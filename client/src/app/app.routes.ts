import { Routes } from '@angular/router';
import { LoginComponent } from './auth/login/login.component';
import { RegisterComponent } from './auth/register/register.component';
import { CallbackComponent } from './auth/callback/callback.component';
import { DefaultPageComponent } from './default-page/default-page.component';
import { authGuard } from './guards/auth.guard';

export const routes: Routes = [
    { path: '', component: DefaultPageComponent, canActivate: [authGuard] },
    { path: 'login', component: LoginComponent },
    { path: 'register', component: RegisterComponent },
    { path: 'auth/callback', component: CallbackComponent }
];
