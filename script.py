import random

with open('users_dump.sql', mode='w') as file:
    file.write("-- PostgreSQL dump for users table\n\n")
    
    for i in range(1, 100001):
        name = f"User{i}"
        roll_no = f"1071{str(i).zfill(5)}"
        email = f"user{i}@nitt.edu"
        pin = str(random.randint(1000, 9999))  
        is_verified = random.choice([True, True])
        is_approved = random.choice([True, True])
        wallet_balance = random.randint(100000, 200000)  
        
        insert_stmt = f"INSERT INTO users (name, roll_no, email, pin, is_verified, is_approved, wallet_balance) VALUES (" \
                      f"'{name}', '{roll_no}', '{email}', '{pin}', {is_verified}, {is_approved}, {wallet_balance});\n"
        
        file.write(insert_stmt)

print("SQL dump file generated successfully.")
