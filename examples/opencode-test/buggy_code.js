// Intentionally buggy JavaScript code for OpenCode to analyze

function calculateAverage(numbers) {
    var total = 0;
    for (var i = 0; i < numbers.length; i++) {
        total += numbers[i];
    }
    return total / numbers.length; // Bug: division by zero if empty array
}

function findMax(arr) {
    var max = 0; // Bug: assumes all numbers are positive
    for (var i = 0; i < arr.length; i++) {
        if (arr[i] > max) {
            max = arr[i];
        }
    }
    return max;
}

function processUserData(users) {
    var results = [];
    for (var i = 0; i < users.length; i++) {
        var user = users[i];
        if (user.age > 18) { // Bug: what if age is undefined?
            results.push({
                name: user.name.toUpperCase(), // Bug: what if name is null?
                email: user.email.toLowerCase(),
                adult: true
            });
        }
    }
    return results;
}

// Example usage with potential issues
var numbers = [];
console.log("Average:", calculateAverage(numbers)); // Will return NaN

var negativeNumbers = [-5, -10, -2];
console.log("Max:", findMax(negativeNumbers)); // Will return 0 instead of -2

var users = [
    { name: "Alice", age: 25, email: "ALICE@EXAMPLE.COM" },
    { name: null, age: 30, email: "BOB@EXAMPLE.COM" }, // Null name
    { name: "Charlie", email: "CHARLIE@EXAMPLE.COM" } // Missing age
];
console.log("Processed:", processUserData(users));