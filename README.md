# goscm
simple scheme interpreter

#substituation rule
<pre>
To evaluate an application:
  Evaluate the operator to get procedure
  Evaluate the operands to get arguments
  Apply the procedure to the arguments
    Copy the body of the procedure.
     subsituting the arguments applied
     for the formal parameters of the pocedure.
    Evaluate the resulting new body.
</pre>
