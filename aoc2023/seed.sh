
path="day${1}"

if [ -d $path ]; then
    echo "already done"
    exit 1
fi

mkdir $path
cp starter/* $path
touch $path/input.txt
