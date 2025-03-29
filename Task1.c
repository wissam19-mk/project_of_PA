#include <stdio.h>
#include <stdlib.h>
#include <string.h>


// constants
#define LIFE 'X'  // ALive cell
#define DEAD '+' // dead cell
#define NEIGHBOR_COUNT 8 // NEIGHBOR_count in function (generate_next_generation)

// readding the first gen form the input file
void create_grid(FILE *input_file, char **first_grid, int N, int M){
  
    char buffer[50];
    fgets(buffer, N+2, input_file); 
    for (int i = 0; i < M; i++) {
         fgets(first_grid[i], N + 2, input_file);  
    }
}

void generate_next_generation(char **first_grid, char **new_gridd, int N, int M){
        
     int i, j; 
  
    // initialize the new generation with default values
    for (i = 0; i < M; i++) {
        for (j = 0; j < N; j++) {
            new_gridd[i][j] = DEAD;
        }
    }

        for (i = 0; i < M; i++){
            for(j = 0; j < N; j++){
                //defualt vlue the nighbor
                char top_left = DEAD;
                char top = DEAD ;
                char top_right = DEAD;
                char left = DEAD;
                char right = DEAD;
                char down_left = DEAD;
                char down = DEAD;
                char down_right = DEAD;
            //check the nighbor if in the bound of matrix 
                if (i > 0 && j > 0)
                    top_left = first_grid[i-1][j-1]; 

                if (i > 0) 
                     top = first_grid[i-1][j];
            
                if (i > 0 && j < N-1)
                     top_right = first_grid[i-1][j+1];
            
                if(j > 0) 
                     left = first_grid[i][j-1];
            
                if (i < M-1 && j > 0) 
                     down_left = first_grid[i+1][j-1];
            
                if (i < M-1) 
                     down = first_grid[i+1][j];
            
                if (i < M-1 && j < N-1) 
                     down_right = first_grid[i+1][j+1];
                if (j < N-1) 
                     right = first_grid[i][j+1];
            
                 //TO SAVE THE VLUE OF EACH NIGHBER
                 const char neighbors[NEIGHBOR_COUNT]= {
                    top_left, top, top_right, left
                    , right , down_left, down, down_right
                };

                // star to check the game conditions if its life do somthing                      
                if (first_grid[i][j] == LIFE){ 
                     int counter = 0 ; 
                    // counter to count how many nigbers the current cell has
                     
                    for(int c = 0; c < 8; c++){
                        if (neighbors[c] == LIFE){
                             counter++; 
                        }
                    }   

                     if(counter < 2){
                    //change the value in the new grid not in the same generation deoend in the condition
                             new_gridd[i][j] = DEAD; 
                    }   
                    else if (counter == 2 || counter == 3) {
                             new_gridd[i][j] = LIFE;
                    }
                    else{
                         new_gridd[i][j] = DEAD;
                    }
                }
                // to count and check the dead nighbors
                else if (first_grid[i][j] == DEAD){
                     int counter = 0 ; 
                     for(int c = 0; c < 8; c++){
                        if (neighbors[c] == LIFE){
                             counter++; 
                        }
                     }
                     if(counter == 3){
                         new_gridd[i][j] = LIFE;
                     }
                }   
            }
        }
        
}


    // copy the new genaration
    void copy_next_gen(char **first_grid, char **New_gridd, int N, int M){
        int i, j; 
        
        for (i = 0; i < M; i++) {
            for (j = 0; j < N; j++) {
                first_grid[i][j] = New_gridd[i][j];
            }
        }
    }

    void copy_in_output_file(FILE *output_file, char** New_gridd , int N, int M){
       
        // wirte in the output file
        for (int i = 0; i < M; i++) {
            for (int j = 0; j < N; j++) {
                fprintf(output_file,"%c",New_gridd[i][j]);
            }
            fprintf(output_file,"\n");
        }
        fprintf(output_file,"\n");
    }


int main (int argc, char **argv){
    
    //variable declare
    FILE *input_file = fopen(argv[1], "r"); //reading only
    FILE *output_file = fopen(argv[2], "w");
    char **first_grid,**New_gridd;
    int T, N, M, K; 
    int generation;

    //test input file 
    if (input_file == NULL){
        fclose(output_file);
        printf("Failed to open input file");  
        return -1;
    }
    //test output file if failed to find then close the input file to save memory
    if (output_file == NULL){
        fclose(input_file);
        printf("Failed to open output file");  
        return -1;
    }
    
    
    // readding from the input file
    fscanf(input_file,"%d",&T);
    fscanf(input_file,"%d %d",&M,&N);
    fscanf(input_file,"%d",&K);
    
    //Allcation in the dainamic as a vector
    first_grid = (char**)malloc(M*(sizeof(char*)));
         
        // to check the vector of pointer if its alocated in the memory
        // if failed to allocate first_grid then close the files
        if (first_grid == NULL) {
            fclose(input_file);
            fclose(output_file);
            printf("Erorr in allocat in dinamic location in current grid");
            return -1;
        }
    // to make matrix of pionter
        for(int i = 0; i < M; i++){
            first_grid[i] = (char*) malloc ((N+2)*(sizeof(char)));
        }

     //Allcation in the dainamic as a vector for the new genaration 
    New_gridd = (char**)malloc(M*(sizeof(char*)));

     // to check the vector of pointer if its alocated in the memory
        if (New_gridd == NULL) {
            fclose(input_file);
            fclose(output_file);
            printf("Erorr in allocat in dinamic location in New grid ");
            return -1;
        }
    // to make matrix of pionter for the new genaration 
        for(int i = 0 ; i < M; i++){
            New_gridd[i] = (char*) malloc ((N+2)*(sizeof(char)));
        }


    //funtion call for the creat the fisrt ganaretion 
    create_grid(input_file, first_grid , N, M);
    copy_in_output_file(output_file ,first_grid , N, M);
    
    for (generation = 0; generation < K; generation++){
            generate_next_generation(first_grid, New_gridd, N, M);
        // copy the next gen to first matrix to contiune creation of new genation  
            copy_next_gen(first_grid, New_gridd, N, M);
        //write the new generation to output file ;
            copy_in_output_file(output_file, New_gridd ,N , M);
    }
   // free the memory  the two matrixs 
    for(int i = 0; i < M; i++){
        free (New_gridd[i]);
    }
    free(New_gridd);
    for(int i = 0; i < M; i++){
       free(first_grid[i]);
    }
    free(first_grid);
    
    //close the input file and output file after uesing
    fclose(input_file);
    fclose(output_file);
    
    return 0;
} 