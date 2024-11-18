
async function loginAndGetSessionToken() {
    const loginPayload = {
        email: 'test@example.com',  // Replace with real credentials
        pin: '$2a$10$fN/kwfGtG.qf52MZIvz13.tx9hm/9DWaAaIwt6ruBG3LZnuV9MgBC'                 // Replace with actual pin
    };
  
    // Send the POST request with login details
    const loginRes = await fetch('http://localhost:7070/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(loginPayload)
    });
  
    // Check if login was successful (status code 200)
    if (loginRes.ok) {
      console.log('Login successful');
  
      // Get the session token from the response cookies (if it exists)
      const cookies = loginRes.headers.get('set-cookie');
      const sessionToken = extractSessionToken(cookies);
  
      if (sessionToken) {
        console.log('Session Token:', sessionToken);
        return sessionToken;
      } else {
        console.error('Session token not found in cookies.');
        return null;
      }
    } else {
      console.error('Login failed with status:', loginRes.status);
      return null;
    }
  }
  
  // Function to extract session token from cookies
  function extractSessionToken(cookieHeader) {
    // Look for a cookie with the name 'session_token' in the set-cookie header
    const tokenRegex = /session_token=([^;]+)/;
    const match = cookieHeader ? cookieHeader.match(tokenRegex) : null;
    return match ? match[1] : null;
  }
  
  // Call the login function and get the token
  loginAndGetSessionToken().then((token) => {
    if (token) {
      console.log('Got session token:', token);
    } else {
      console.log('Failed to retrieve session token.');
    }
  });
  