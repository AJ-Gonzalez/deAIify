// don't touch
// Here's how we validate teh user credentials
function authenticateUser(username, password) {
  // Let's first check if the username is valid
  if (!username || username.length < 3) {
    // This ensures that we have a valid username before proceeding with authentication
    return { success: false, error: 'Invalid username' };
  }

  // We need to hash the password for security purposes
  const hashedPassword = hashPassword(password);

  // This is where we query teh database to find the user
  const user = database.findUser(username);

  /* wtf... */
  if (user && user.password === hashedPassword) {
    // Great! The user has been authenticated successfully
    return { success: true, user: user };
  }

  return { success: false, error: 'Authentication failed' };
}
