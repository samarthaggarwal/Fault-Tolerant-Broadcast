from matplotlib import pyplot as plt
import pandas as pd
import math

pd.set_option("display.max_rows", None, "display.max_columns", None)
pd.set_option('display.width', 150)

file = "data/vm.csv"
df0 = pd.read_csv(file)
df0["DSTotalMsgCount"] = df0.numNodes * (df0.numNodes-1)
df0["DSHonestMsgCount"] = (df0.numNodes-df0.realFaults) * (df0.numNodes-1)
df0["DSFailedMsgCount"] = df0.realFaults * (df0.numNodes-1)

# Graph1: Vary n with constant f
faults = [int(df0.realFaults.quantile(i/10)) for i in range(11)]
faults = list(set(faults))
faults.sort()
for f in faults:
	df = df0[ df0.realFaults==f ]
	df = df.groupby('numNodes').mean().reset_index()
	df["ratio"] = df.totalMsgCount / df.DSTotalMsgCount

	#plt.figure()
	plt.plot(df.numNodes, df.ratio)
	plt.title(f"Ratio of Total Messages for {f} Faults")
	plt.xlabel("Number of Nodes")
	plt.ylabel("Ratio of Messages in Galil-Mayer to Dolev-Strong")
	plt.tight_layout()
	plt.savefig(f"plots/f{f}.png", pad_inches=4)
	plt.clf()

# Graph2: Vary f with constant n
for n in list(df0.numNodes.unique()):
	df = df0[ df0.numNodes==n ]
	df = df.groupby('realFaults').mean().reset_index()
	df["additionalMsg"] = df.totalMsgCount - df.totalMsgCount[ df.realFaults==0 ][0]
	df["ratio"] = df.additionalMsg / math.sqrt(df.totalMsgCount[ df.realFaults==0 ][0])
	
	#plt.figure()
	plt.plot(df.realFaults, df.ratio, label="Empirical Value")
	plt.plot(df.realFaults, df.realFaults, label="Theoretical Bound")
	plt.title(f"Ratio of Additional Messages due to Faults to sqrt(n) for {n} Nodes")
	plt.xlabel("Number of Faults")
	plt.ylabel("Ratio of Additional Messages to sqrt(n)")
	plt.legend()
	plt.tight_layout()
	plt.savefig(f"plots/n{n}.png", pad_inches=4)
	plt.clf()

