from matplotlib import pyplot as plt
import pandas as pd
import math

pd.set_option("display.max_rows", None, "display.max_columns", None)
pd.set_option('display.width', 150)

file = "data/worst_case.csv"
df0 = pd.read_csv(file)
df0["DSTotalMsgCount"] = df0.numNodes * (df0.numNodes-1)
df0["DSHonestMsgCount"] = (df0.numNodes-df0.realFaults) * (df0.numNodes-1)
df0["DSFailedMsgCount"] = df0.realFaults * (df0.numNodes-1)

# Graph3: Vary n, f=n-1
#df = df0[ df0.numNodes==n ]
df = df0.groupby('numNodes').mean().reset_index()
df["theoreticalBound"] = df.numNodes + (df.realFaults*(df.numNodes**0.5))
#df["ratio"] = df.totalMsgCount / df.theoreticalBound
df["lowerBound"] = df.numNodes*((df.numNodes**0.5) + 1)/2 - 1

#plt.figure()
plt.plot(df.numNodes, df.theoreticalBound, label="Theoretical Upper Bound")
plt.plot(df.numNodes, df.totalMsgCount, label="Empirical Value")
plt.plot(df.numNodes, df.lowerBound, label="Derived Lower Bound")
plt.title(f"Message Count in Potential Worst-Case")
plt.xlabel("Number of Nodes")
plt.ylabel("Total Messages")
plt.legend()
plt.tight_layout()
plt.savefig(f"plots/worstcase.png", pad_inches=4)
plt.clf()

