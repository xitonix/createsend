## General Guidelines for Contributing
**First thing first! Let's offload some chores**

Make sure you setup [pre-commit](https://pre-commit.com/) on your machine. This will keep your books tidy and runs a bunch of linting and formatting tasks before each commit.

**Size Matters**

If you you are adding more CM API functions to the wrapper, make sure to add a few functions at a time. This will help us to focus on a small area of work and test it exhaustively as we move forward. 

**Test test test**

- We should keep the unit test coverage at 100% at all times!
- Create a small playground on your machine and use the wrapper to **manually test** the new code. We are building **for developers**, let's test it **as developers**. 

**Errors**

- All the top level wrapper methods must return a [createsend error](https://github.com/xitonix/createsend/blob/master/errors.go) in case they are returning a custom error. Keep in mind that all the methods of `internal.Client` already return `createsend errors`, so you don't need to wrap them.

    Example

      func (a *clientsAPI) SuppressionList(clientID string,
        pageSize, page int,
        orderBy order.SuppressionListField,
        direction order.Direction) (*clients.SuppressionList, error) {

        result := new(internal.SuppressionList)
        err := a.client.Get(path, &result)
        if err != nil {
          return nil, err <==== No Need to Wrap this error
        }

        list, err := result.ToSuppressionList()
        if err != nil {
          return nil, newClientError(ErrCodeDataProcessing) <==== HERE
        }

        return list, nil
      }
  
- All the top level wrapper methods must explicitly return the empty value of the desired response (or `nil`) on error:

     Example

      // YES
      result := new(internal.SuppressionList)
      err := a.client.Get(path, &result)
      if err != nil {
        return nil, err
      }

      // NO (Even if the result is the same)
      result := new(internal.SuppressionList)
      err := a.client.Get(path, &result)
      return result, err


**Documentation**

All the exported types, fields and methods **MUST** be documented. That said, we should avoid documenting CM API's behaviour as much as possible. How much is too much? We should use common sense.

Example:

```
// PerformAction performs XYZ on ABC.
//
// howManyTimes must be between 1 and 5 <== This line must be avoided
func PerformAction(howManyTimes int) error {
}
```
