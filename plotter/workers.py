#!/usr/bin/env python3

from matplotlib import pyplot as plt
import numpy as np
import csv
import pandas as pd

n_jobs=10
jobs=['lint', 'test', 'sonar', 'build', 'cont', 'depend', 'pre', 'bench', 'verify', 'deploy']
#jobs=[1, 2, 3, 4, 5, 6, 7, 8, 9]
job_col_index=[4,5,6,8,9,10,12,13,14,16]  # column indices of jobs

def readFile(name):
    #code_lines=[]
    n_workers=[]
    job_lat=[]
    q_lat=[]

    # read lines and number of test rounds
    with open(name) as file:
        reader = csv.reader(file)
        row_n = 0
        for row in reader:
            if row_n > 6:
                if row[0] != '':
                    n_workers.append(int(row[0]))
            row_n += 1
    
    # create 3D matrix
    for i in range(len(n_workers)):
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
            if row_n > 6:
                if row[0] != '':
                    test_case += 1
                for j in range(n_jobs):
                    if row[3] == 'job':
                        job_lat[test_case][j][test_attempt] = float(row[job_col_index[j]])
                    if row[3] == 'queued':
                        q_lat[test_case][j][test_attempt] = float(row[job_col_index[j]])
                        
                if row[3] == 'queued':
                    if test_attempt == 2:
                        test_attempt = 0
                    else:
                        test_attempt += 1
            row_n += 1
    
    return n_workers, job_lat, q_lat

def parseData(n_workers, job_lat, q_lat):
    avg_job_lat = []
    err_job_lat = []
    avg_q_lat = []
    err_q_lat = []
    
    for i in range(len(n_workers)):
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

def aggrLat(n_workers, avg_job_lat, avg_q_lat):
    aggr_lat = []
    for i in range(len(n_workers)):
        total_lat = 0.0
        aggr_lat.append([])
        for j in range(n_jobs):
            total_lat += avg_job_lat[i][j] + avg_q_lat[i][j]
            aggr_lat[i].append(total_lat)
    
    return aggr_lat
            
def plotJobLineGraph(aggr_job_lat, err_job_lat):
    fig, ax = plt.subplots()
    ax.errorbar(jobs, aggr_job_lat[0], yerr=err_job_lat[0], ecolor="red", color="green", label="concurrency=1")
    ax.errorbar(jobs, aggr_job_lat[1], yerr=err_job_lat[1], ecolor="red", color="cornflowerblue", label="concurrency=2")
    #ax.set_xlabel('job')
    ax.set_ylabel('cumulative average latency (s)')
    #ax.set_title('Latency of CICD pipelines for different number of workers')
    plt.rc('grid', linestyle="--", color='#C6C6C6')
    #plt.xticks(rotation=90)
    plt.legend()
    plt.grid()
    plt.savefig('docs/imgs/total_latency_workers.pdf', bbox_inches="tight")
    plt.show()
    
def plotHistogram(typ, avg_lat, err_lat):
    clrs={'concurrency=1': 'green', 'concurrency=2': 'cornflowerblue'}
    low_vals = []
    high_vals = []
    labels = []    
    low_err = []
    high_err = []
    
    for i in range(len(jobs)):
        if typ == 'job' and i == 7:
            continue
        else:
            low_vals.append(avg_lat[0][i])
            high_vals.append(avg_lat[1][i])
            labels.append(jobs[i])            
            low_err.append(err_lat[0][i])
            high_err.append(err_lat[1][i])

    df = pd.DataFrame({'concurrency=1': low_vals, 'concurrency=2': high_vals}, index=labels)
    ax = df.plot.bar(rot=0, yerr={'concurrency=1': low_err, 'concurrency=2': high_err}, ecolor="crimson", color=clrs)
    ax.set_ylabel('average latency (s)')
    if typ == 'job':
        #ax.set_title('Average execution latency of jobs for different number of workers')
        plt.rc('grid', linestyle="--", color='#C6C6C6')
        plt.grid()
        #plt.xticks(rotation=90)
        plt.savefig('docs/imgs/job_execution_histogram_workers.pdf', bbox_inches="tight")
        plt.show()
    else:
        #ax.set_title('Average queued latency of jobs for different number of workers')
        plt.rc('grid', linestyle="--", color='#C6C6C6')
        plt.grid()
        #plt.xticks(rotation=90)
        plt.savefig('docs/imgs/job_queued_histogram_workers.pdf', bbox_inches="tight")
        plt.show()

n_workers, job_lat, q_lat = readFile('docs/results/workers-vs-pipelines.csv')
avg_job_lat, avg_q_lat, err_job_lat, err_q_lat = parseData(n_workers, job_lat, q_lat)
aggr_job_lat = aggrLat(n_workers, avg_job_lat, avg_q_lat)
plotJobLineGraph(aggr_job_lat, err_job_lat)

plotHistogram('job', avg_job_lat, err_job_lat)
plotHistogram('queued', avg_q_lat, err_q_lat)



