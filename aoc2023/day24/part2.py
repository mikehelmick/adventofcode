from z3 import Int, Solver, sat


solver = Solver()
x, y, z, vx, vy, vz = map(Int, ('x', 'y', 'z', 'vx', 'vy', 'vz'))

t1 = Int("t1")
solver.add(t1 > 0)
solver.add(x + vx * t1 == 237822270988608 + 115 * t1)
solver.add(y + vy * t1 == 164539183264530 + 346 * t1)
solver.add(z + vz * t1 == 381578606559948 + -342 * t1)

t2 = Int("t2")
solver.add(t2 > 0)
solver.add(x + vx * t2 == 287838354624648 + -5 * t2)
solver.add(y + vy * t2 == 284335343503076 + -84 * t2)
solver.add(z + vz * t2 == 181128681512377 + 175 * t2)

t3 = Int("t3")
solver.add(t3 > 0)
solver.add(x + vx * t3 == 341046208911993 + -74 * t3)
solver.add(y + vy * t3 == 120694764237967 + 129 * t3)
solver.add(z + vz * t3 == 376069872241870 + -78 * t3)


if solver.check() == sat:
  m = solver.model()
  aX = m.eval(x).as_long()
  aY = m.eval(y).as_long()
  aZ = m.eval(z).as_long()
  ans = aX + aY + aZ
  print(ans)
  print('answers {aX}, {aY}, {aZ} == {ans}')
else:
  print("unsat")

