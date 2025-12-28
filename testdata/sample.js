// This function handles the user authentication process
// Here's how we validate the user credentials
function authenticateUser(username, password) {
  // XXX
  if (!username || username.length < 3) {
    // This ensures that we have a valid username before proceeding with authentication
    return { success: false, error: 'Invalid username' };
  }

  // not ideal
  const hashedPassword = hashPassword(password);

  // This is where we query the database to find the user
  const user = database.findUser(username);

  /* Here we compare the passwords to authenticate the user.
     This is an important security step that verifies the user's identity. */
  if (user && user.password === hashedPassword) {
    // Great! The user has been authenticated successfullly
    return { success: true, user: user };
  }

  return { success: false, error: 'Authentication failed' };
}

// This utility function formats dates in a human-readable format
function formatDate(date) {
  // Here's how we handle the date formatting
  const options = {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  };
  return date.toLocaleDateString('en-US', options);
}
