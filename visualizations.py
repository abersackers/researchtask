# CSV analysis in python
# Using pandas and numpy, takes in a csv file that is outputted from getinfo.go
# and prints average bytes per get request as well as calculating the cdf
# and printing some percentiles of interests

import pandas as pd
import numpy as np
import matplotlib.pyplot as plt



# takes in a pandas dataframe with the assumption that the 5th indexed
# entry is the total bytes of the get request and calculates an average
# of the total bytes of every entry
def calculate_average_bytes(df):
    total_entries = 0
    total_bytes = 0
    for index, row in df.iterrows():
        total_entries += 1
        total_bytes += row[5]

    return 1.0 * total_bytes / total_entries


# Calculate cdf using given function from Sudhessh
def cdf(data):
    n = len(data)
    x = np.sort(data)  # sort your data
    y = np.arange(1, n + 1) / n  # calculate cumulative probability
    return x, y


# Using an array of percentiles of interest, print out the associated value
# of that percentile from an inputted cdf
def print_percentiles(sorted):
    percentiles = [50, 70, 75, 80, 85, 90, 95, 99]
    for i in percentiles:
        print(str(i) + "th percentile is: " + str(np.percentile(sorted, i)))


def line_plot(start_times, total_times):
    plt.style.use('seaborn-whitegrid')
    plt.plot(start_times, total_times)
    plt.savefig('start_v_total.png')


# read the csv output of getinfo.go
data = pd.read_csv("requestData.csv")

# calculate & print percentile data for difference in time
start_times = pd.array(data['Start Time']).to_numpy()
total_times = pd.array(data['End Time'] - data['Start Time']).to_numpy()
sorted, cdf = cdf(total_times)
print_percentiles(sorted)
line_plot(start_times, total_times)
