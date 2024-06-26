#!/usr/bin/env python3

from matplotlib import pyplot as plt
import numpy as np
import csv
import pandas as pd

n_jobs=10
jobs=['lint', 'test', 'sonar', 'build', 'cont', 'depend', 'pre', 'bench', 'verify', 'deploy']
#jobs=[1, 2, 3, 4, 5, 6, 7, 8, 9]
job_col_index=[5,6,7,9,10,11,13,14,15,17]  # column indices of jobs

def readFile(name):
    code_lines=[]
    test_lines=[]
    total_lines=[]

    job_lat=[]
    q_lat=[]

    # read lines and number of test rounds
    with open(name) as file:
        reader = csv.reader(file)
        row_n = 0
        for row in reader:
            if row_n > 8:
                if row[0] != '':
                    code_lines.append(int(row[0]))
                    test_lines.append(int(row[1]))
                    total_lines.append(int(row[2]))
            row_n += 1

    # create 3D matrix
    for i in range(len(total_lines)):
        job_lat.append([])
        q_lat.append([])
        for j in range(n_jobs):
            job_lat[i].append([0, 0, 0])
            q_lat[i].append([0, 0, 0])

    # add latency values
    with open(name) as file:
        reader = csv.reader(file)
        test_case = -1 # updated to number of test cases when reading the file
        test_attempt = 0
        row_n = 0
        for row in reader:
            # skip header
            if row_n > 8:
                if row[0] != '':
                    test_case += 1
                for j in range(n_jobs):
                    if row[4] == 'job':
                        job_lat[test_case][j][test_attempt] = float(row[job_col_index[j]])
                    if row[4] == 'queued':
                        q_lat[test_case][j][test_attempt] = float(row[job_col_index[j]])

                if row[4] == 'queued':
                    if test_attempt == 2:
                        test_attempt = 0
                    else:
                        test_attempt += 1
            row_n += 1

    return total_lines, job_lat, q_lat

def parseData(total_lines, job_lat, q_lat):
    avg_job_lat = []
    err_job_lat = []
    avg_q_lat = []
    err_q_lat = []

    for i in range(len(total_lines)):
        avg_job_lat.append([])
        err_job_lat.append([])
        avg_q_lat.append([])
        err_q_lat.append([])
        for j in range(n_jobs):
            avg_job_lat[i].append(np.mean(job_lat[i][j]))
            err_job_lat[i].append(np.std(job_lat[i][j]))
            avg_q_lat[i].append(np.mean(q_lat[i][j]))
            err_q_lat[i].append(np.std(q_lat[i][j]))

    return avg_job_lat, avg_q_lat, err_job_lat, err_q_lat

def aggrLat(total_lines, avg_job_lat, avg_q_lat):
    aggr_lat = []
    for i in range(len(total_lines)):
        total_lat = 0.0
        aggr_lat.append([])
        for j in range(n_jobs):
            total_lat += avg_job_lat[i][j] + avg_q_lat[i][j]
            aggr_lat[i].append(total_lat)

    return aggr_lat

def plotJobLineGraph(aggr_job_lat, err_job_lat):    
    fig, ax = plt.subplots()
    ax.errorbar(jobs, aggr_job_lat[2], yerr=err_job_lat[2], ecolor="red", color="goldenrod", label="high")
    ax.errorbar(jobs, aggr_job_lat[1], yerr=err_job_lat[1], ecolor="red", color="cornflowerblue", label="mid")
    ax.errorbar(jobs, aggr_job_lat[0], yerr=err_job_lat[0], ecolor="red", color="green", label="low")
    #ax.set_xlabel('job')
    ax.set_ylabel('cumulative average latency (s)')
    #ax.set_title('Total latency of CICD pipelines for different chaincodes')
    plt.rc('grid', linestyle="--", color='#C6C6C6')
    #plt.xticks(rotation=90)
    plt.legend()
    plt.grid()
    plt.savefig('docs/imgs/total_latency_pipelines.pdf', bbox_inches="tight")
    plt.show()

def plotHistogram(typ, avg_lat, err_lat):
    clrs={'low': 'green', 'mid': 'cornflowerblue', 'high': 'goldenrod'}
    low_vals = []
    mid_vals = []
    high_vals = []
    labels = []

    low_err = []
    mid_err = []
    high_err = []

    for i in range(len(jobs)):
        if typ == 'job' and i == 7:
            continue
        else:
            low_vals.append(avg_lat[0][i])
            mid_vals.append(avg_lat[1][i])
            high_vals.append(avg_lat[2][i])
            labels.append(jobs[i])

            low_err.append(err_lat[0][i])
            mid_err.append(err_lat[1][i])
            high_err.append(err_lat[2][i])

    df = pd.DataFrame({'low': low_vals, 'mid': mid_vals, 'high': high_vals}, index=labels)
    ax = df.plot.bar(rot=0, yerr={'low': low_err, 'mid': mid_err, 'high': high_err}, ecolor="crimson", color=clrs)
    ax.set_ylabel('average latency (s)')
    if typ == 'job':
        #ax.set_title('Average execution latency of jobs for different chaincodes')
        plt.rc('grid', linestyle="--", color='#C6C6C6')
        plt.grid()
        #plt.xticks(rotation=90)
        plt.savefig('docs/imgs/job_execution_histogram.pdf', bbox_inches="tight")
        plt.show()
    else:
        #ax.set_title('Average queued latency of jobs for different chaincodes')
        plt.rc('grid', linestyle="--", color='#C6C6C6')
        plt.grid()
        #plt.xticks(rotation=90)
        plt.savefig('docs/imgs/job_queued_histogram.pdf', bbox_inches="tight")
        plt.show()

total_lines, job_lat, q_lat = readFile('docs/results/code-vs-pipelines-extended.csv')
avg_job_lat, avg_q_lat, err_job_lat, err_q_lat = parseData(total_lines, job_lat, q_lat)
aggr_lat = aggrLat(total_lines, avg_job_lat, avg_q_lat)
plotJobLineGraph(aggr_lat, err_job_lat)

plotHistogram('job', avg_job_lat, err_job_lat)
plotHistogram('queued', avg_q_lat, err_q_lat)

#---------- extreme case-----------------------#

def readNewFile(name):
    code_lines=[]
    test_lines=[]
    total_lines=[]

    job_lat=[]
    q_lat=[]

    # read lines and number of test rounds
    with open(name) as file:
        reader = csv.reader(file)
        row_n = 0
        for row in reader:
            if row_n > 8 and row_n < 33:
                if row[0] != '':
                    code_lines.append(int(row[0]))
                    test_lines.append(int(row[1]))
                    total_lines.append(int(row[2]))
            row_n += 1

    # create 3D matrix
    for i in range(len(total_lines)):
        job_lat.append([])
        q_lat.append([])
        for j in range(n_jobs):
            job_lat[i].append([0, 0, 0])
            q_lat[i].append([0, 0, 0])

    # add latency values
    with open(name) as file:
        reader = csv.reader(file)
        test_case = -1 # updated to number of test cases when reading the file
        test_attempt = 0
        row_n = 0
        for row in reader:
            # skip header
            if row_n > 8 and row_n < 33:
                if row[0] != '':
                    test_case += 1
                for j in range(n_jobs):
                    if row[4] == 'job':
                        job_lat[test_case][j][test_attempt] = float(row[job_col_index[j]])
                    if row[4] == 'queued':
                        q_lat[test_case][j][test_attempt] = float(row[job_col_index[j]])

                if row[4] == 'queued':
                    if test_attempt == 2:
                        test_attempt = 0
                    else:
                        test_attempt += 1
            row_n += 1

    return total_lines, job_lat, q_lat

def parseDataNew(total_lines, job_lat, q_lat):
    avg_job_lat = []
    err_job_lat = []
    avg_q_lat = []
    err_q_lat = []

    for i in range(len(total_lines)):
        avg_job_lat.append([])
        err_job_lat.append([])
        avg_q_lat.append([])
        err_q_lat.append([])
        for j in range(n_jobs):
            avg_job_lat[i].append(np.mean(job_lat[i][j]))
            err_job_lat[i].append(np.std(job_lat[i][j]))
            avg_q_lat[i].append(np.mean(q_lat[i][j]))
            err_q_lat[i].append(np.std(q_lat[i][j]))

    return avg_job_lat, avg_q_lat, err_job_lat, err_q_lat

def plotHistogramNew(typ, avg_lat, err_lat):
    clrs={'low': 'green', 'mid': 'cornflowerblue', 'high': 'goldenrod', 'extreme': 'slategrey'}
    low_vals = []
    mid_vals = []
    high_vals = []
    extreme_vals = []
    labels = []

    low_err = []
    mid_err = []
    high_err = []
    extreme_err = []

    for i in range(len(jobs)):
        if typ == 'job' and i == 7:
            continue
        else:
            low_vals.append(avg_lat[0][i])
            mid_vals.append(avg_lat[1][i])
            high_vals.append(avg_lat[2][i])
            extreme_vals.append(avg_lat[3][i])
            labels.append(jobs[i])

            low_err.append(err_lat[0][i])
            mid_err.append(err_lat[1][i])
            high_err.append(err_lat[2][i])
            extreme_err.append(err_lat[3][i])

    df = pd.DataFrame({'low': low_vals, 'mid': mid_vals, 'high': high_vals, 'extreme': extreme_vals}, index=labels)
    ax = df.plot.bar(rot=0, yerr={'low': low_err, 'mid': mid_err, 'high': high_err, 'extreme': extreme_err}, ecolor="crimson", color=clrs)
    ax.set_ylabel('average latency (s)')
    if typ == 'job':
        #ax.set_title('Average execution latency of jobs for different chaincodes')
        plt.rc('grid', linestyle="--", color='#C6C6C6')
        plt.grid()
        #plt.xticks(rotation=90)
        plt.savefig('docs/imgs/job_execution_extended_histogram.pdf', bbox_inches="tight")
        plt.show()
    else:
        #ax.set_title('Average queued latency of jobs for different chaincodes')
        plt.rc('grid', linestyle="--", color='#C6C6C6')
        plt.grid()
        #plt.xticks(rotation=90)
        plt.savefig('docs/imgs/job_queued_extended_histogram.pdf', bbox_inches="tight")
        plt.show()

total_lines, job_lat, q_lat = readNewFile('docs/results/code-vs-pipelines-extended.csv')
avg_job_lat, avg_q_lat, err_job_lat, err_q_lat = parseDataNew(total_lines, job_lat, q_lat)

plotHistogramNew('job', avg_job_lat, err_job_lat)
plotHistogramNew('queued', avg_q_lat, err_q_lat)
