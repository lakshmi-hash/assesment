import { ComponentFixture, TestBed } from '@angular/core/testing';
import { UserListComponent } from './user-list.component';
import { UserService } from '../../services/user.service';
import { MatDialog } from '@angular/material/dialog';
import { of } from 'rxjs';

describe('UserListComponent', () => {
  let component: UserListComponent;
  let fixture: ComponentFixture<UserListComponent>;
  let userServiceSpy: jasmine.SpyObj<UserService>;

  beforeEach(async () => {
    const spy = jasmine.createSpyObj('UserService', ['getUsers', 'deleteUser']);
    const dialogSpy = jasmine.createSpyObj('MatDialog', ['open']);

    await TestBed.configureTestingModule({
      declarations: [UserListComponent],
      providers: [
        { provide: UserService, useValue: spy },
        { provide: MatDialog, useValue: dialogSpy }
      ]
    }).compileComponents();

    userServiceSpy = TestBed.inject(UserService) as jasmine.SpyObj<UserService>;
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(UserListComponent);
    component = fixture.componentInstance;
    userServiceSpy.getUsers.and.returnValue(of([{ id: 1, name: 'Test', email: 'test@example.com' }]));
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should load users on init', () => {
    expect(component.users.length).toBe(1);
    expect(component.users[0].name).toBe('Test');
  });
});