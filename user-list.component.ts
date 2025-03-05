import { Component, OnInit } from '@angular/core';
import { UserService } from '../../services/user.service';
import { MatDialog } from '@angular/material/dialog';
import { UserFormComponent } from '../user-form/user-form.component';

@Component({
  selector: 'app-user-list',
  templateUrl: './user-list.component.html',
  styleUrls: ['./user-list.component.css']
})
export class UserListComponent implements OnInit {
  displayedColumns = ['name', 'email', 'actions'];
  users: any[] = [];

  constructor(private userService: UserService, public dialog: MatDialog) {}

  ngOnInit() {
    this.loadUsers();
  }

  loadUsers() {
    this.userService.getUsers().subscribe(users => this.users = users);
  }

  openCreateDialog() {
    const dialogRef = this.dialog.open(UserFormComponent, { data: {} });
    dialogRef.afterClosed().subscribe(result => {
      if (result) this.loadUsers();
    });
  }

  openEditDialog(user: any) {
    const dialogRef = this.dialog.open(UserFormComponent, { data: user });
    dialogRef.afterClosed().subscribe(result => {
      if (result) this.loadUsers();
    });
  }

  deleteUser(id: number) {
    this.userService.deleteUser(id).subscribe(() => this.loadUsers());
  }
}