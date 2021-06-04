MERGED="downloads/merged"
FOLDERS=$(ls ./downloads/bugwood)
# Delete the merged folder before writing toit agin.
rm -rf $MERGED
SITES=$(ls -d ./downloads/*/)
echo $SITES

mkdir downloads/merged
for tree in $FOLDERS
do
    COUNTER=0
    mkdir downloads/merged/$tree
    for site in $SITES
    do
            test=$(ls -l $site$tree)
            if [ $? -ne 0 ]; then
                    echo $tree "was not scraped in " $site
            else
                echo "Copying from $tree from $site"
                for filename in $site$tree/*.jpg
                do
                    echo $COUNTER
                    cp $filename downloads/merged/$tree/$COUNTER.jpg
                    ((COUNTER++))
                done
            fi

                    
    done
done
