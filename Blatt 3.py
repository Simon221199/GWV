# -*- coding: utf-8 -*-
"""
Created on Wed Oct 30 18:14:20 2019

@author: simon
"""


#Reading environment from text file
#Than putting Input into String with removing /n
env = open("blatt3_environment.txt", "r")
envs = "" 
for line in env:
    envs += line.rstrip()

#Converting string from text file to 2D array
w, h, w_, h_ = 20, 10, 0, 0;
env_matrix = [[0 for x in range(w)] for y in range(h)]
index = 0
while index < len(envs):
  letter = envs[index]
  env_matrix[h_][w_] = letter
  w_ += 1
  if w_ == 20:
      w_ = 0
      h_ += 1
  index = index + 1           
  
#finished
env.close()
print(env_matrix)
        
#prints the current search state
def printSearchState():
    positionOfS = [0,0]
    for i in range(len(env_matrix)):
        for j in range(len(env_matrix[i])):
            if env_matrix[i][j] == 's':
                positionOfS = [i,j]
    print(positionOfS)
    positionOfG = [0,0]
    for i in range(len(env_matrix)):
        for j in range(len(env_matrix[i])):
            if env_matrix[i][j] == 'g':
                positionOfG = [i,j]
    print(positionOfS, positionOfG)
    return positionOfS

#Help function to get the goal field's coordinates
def getGoalField():
    positionOfG = [0,0]
    for i in range(len(env_matrix)):
        for j in range(len(env_matrix[i])):
            if env_matrix[i][j] == 'g':
                positionOfG = [i,j]
    return positionOfG
    
#Moves the robot in directions up, down, left or right
#default is none, incorrect input also results to none
#In this presenting way, the robot's position in marked
#with s, so the s is moving from the starting field
def move(direction = 'None'):
    por = printSearchState() # = position of robot
    goalField = getGoalField()
    #Moving the Robot up (if possible) means to swap the field with 's' with the field with ' ' above it
    #Equal for the other directions
    if(direction == 'up'):
        if por[0] > 0 and env_matrix[por[0]-1][por[1]] != 'x':
            env_matrix[por[0]-1][por[1]] = 's'
            env_matrix[por[0]][por[1]] = ' '
    if(direction == 'down'):
        if por[0] < 9 and env_matrix[por[0]+1][por[1]] != 'x':
            env_matrix[por[0]+1][por[1]] = 's'
            env_matrix[por[0]][por[1]] = ' '
    if(direction == 'left'):
        if por[1] > 0 and env_matrix[por[0]][por[1]-1] != 'x':
            env_matrix[por[0]][por[1]-1] = 's'
            env_matrix[por[0]][por[1]] = ' '
    if(direction == 'right'):
        if por[1] < 19 and env_matrix[por[0]][por[1]+1] != 'x':
            env_matrix[por[0]][por[1]+1] = 's'
            env_matrix[por[0]][por[1]] = ' '
            
    #Checking, if goal state is reached (s = g)
    if (por == goalField()):
        print('Success: Goal field reached!')
    print(env_matrix)
    
    
    
    
    