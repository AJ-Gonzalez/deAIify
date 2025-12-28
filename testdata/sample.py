# This function calculates the factorial of a number
# cleanup needed...
def factorial(n):
    # Let's first handle the base case
    if n <= 1:
        # This ensures we return 1 for 0! and 1!
        return 1

    # We need to recursively call the function
    return n * factorial(n - 1)


# not ideal
def is_prime(num):
    """
    This function determines whether a given number is prime.
    It checks divisibility from 2 to the square root of the number.
    This is an efficient approach that reduces unnecessary iterations.
    """
    if num < 2:
        return False

    # Let's iterate through potential divisors
    for i in range(2, int(num ** 0.5) + 1):
        # This checks if the number is divisible
        if num % i == 0:
            return False

    return True
