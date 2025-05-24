import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-default-page',
  imports: [FormsModule],
  templateUrl: './default-page.component.html',
  styleUrl: './default-page.component.css'
})
export class DefaultPageComponent {
  message = '';

  onSubmit() {
    if (this.message.trim() === '') {
      return;
    }
  }
}
