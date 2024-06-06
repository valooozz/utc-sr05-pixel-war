mkfifo /tmp/in_A1 /tmp/out_A1
mkfifo /tmp/in_C1 /tmp/out_C1
mkfifo /tmp/in_N1 /tmp/out_N1

mkfifo /tmp/in_A2 /tmp/out_A2
mkfifo /tmp/in_C2 /tmp/out_C2
mkfifo /tmp/in_N2 /tmp/out_N2

mkfifo /tmp/in_A3 /tmp/out_A3
mkfifo /tmp/in_C3 /tmp/out_C3
mkfifo /tmp/in_N3 /tmp/out_N3

mkfifo /tmp/in_A4 /tmp/out_A4
mkfifo /tmp/in_C4 /tmp/out_C4
mkfifo /tmp/in_N4 /tmp/out_N4

mkfifo /tmp/in_A5 /tmp/out_A5
mkfifo /tmp/in_C5 /tmp/out_C5
mkfifo /tmp/in_N5 /tmp/out_N5

mkfifo /tmp/in_A6 /tmp/out_A6
mkfifo /tmp/in_C6 /tmp/out_C6
mkfifo /tmp/in_N6 /tmp/out_N6

#SITE 1
go run app-base -n A1 < /tmp/in_A1 > /tmp/out_A1 &
go run app-control -n C1 -nbsites 6 < /tmp/in_C1 > /tmp/out_C1 &
go run app-net -n N1 -r '[3,2]' -nbsites 6 -port 4445 < /tmp/in_N1 > /tmp/out_N1 &

#SITE 2
go run app-base -n A2 < /tmp/in_A2 > /tmp/out_A2 &
go run app-control -n C2 -nbsites 6 < /tmp/in_C2 > /tmp/out_C2 &
go run app-net -n N2 -r '[1,6;6,3]' -nbsites 6 -port 4447 < /tmp/in_N2 > /tmp/out_N2 &

#SITE 3
go run app-base -n A3 < /tmp/in_A3 > /tmp/out_A3 &
go run app-control -n C3 -nbsites 6 < /tmp/in_C3 > /tmp/out_C3 &
go run app-net -n N3 -r '[2,4;4,1]' -nbsites 6 -port 4449 < /tmp/in_N3 > /tmp/out_N3 &

#SITE 4
go run app-base -n A4 < /tmp/in_A4 > /tmp/out_A4 &
go run app-control -n C4 -nbsites 6 < /tmp/in_C4 > /tmp/out_C4 &
go run app-net -n N4 -r '[3,5;5,3]' -nbsites 6 -port 4451 < /tmp/in_N4 > /tmp/out_N4 &

#SITE 5
go run app-base -n A5 < /tmp/in_A5 > /tmp/out_A5 &
go run app-control -n C5 -nbsites 6 < /tmp/in_C5 > /tmp/out_C5 &
go run app-net -n N5 -r '[4,4]' -nbsites 6 -port 4453 < /tmp/in_N5 > /tmp/out_N5 &


#SITE 6
go run app-base -n A6 < /tmp/in_A6 > /tmp/out_A6 &
go run app-control -n C6 -nbsites 6 < /tmp/in_C6 > /tmp/out_C6 &
go run app-net -n N6 -r '[2,2]' -nbsites 6 -port 4455 < /tmp/in_N6 > /tmp/out_N6 &


cat /tmp/out_A1 > /tmp/in_C1 &
cat /tmp/out_C1 | tee /tmp/in_A1 > /tmp/in_N1 &
#CONNEXIONS INTERSITES
cat /tmp/out_N1 | tee /tmp/in_C1 | tee /tmp/in_N2 > /tmp/in_N3 &

cat /tmp/out_A2 > /tmp/in_C2 &
cat /tmp/out_C2 | tee /tmp/in_A2 > /tmp/in_N2 &
#CONNEXIONS INTERSITES
cat /tmp/out_N2 | tee /tmp/in_C2 | tee /tmp/in_N3 | tee /tmp/in_N6 > /tmp/in_N1 &

cat /tmp/out_A3 > /tmp/in_C3 &
cat /tmp/out_C3 | tee /tmp/in_A3 > /tmp/in_N3 &
#CONNEXIONS INTERSITES
cat /tmp/out_N3 | tee /tmp/in_C3 | tee /tmp/in_N1 | tee /tmp/in_N2 > /tmp/in_N4 &

cat /tmp/out_A4 > /tmp/in_C4 &
cat /tmp/out_C4 | tee /tmp/in_A4 > /tmp/in_N4 &
#CONNEXIONS INTERSITES
cat /tmp/out_N4 | tee /tmp/in_C4 | tee /tmp/in_N3 > /tmp/in_N5 &

cat /tmp/out_A5 > /tmp/in_C5 &
cat /tmp/out_C5 | tee /tmp/in_A5 > /tmp/in_N5 &
#CONNEXIONS INTERSITES
cat /tmp/out_N5 | tee /tmp/in_C5 | tee /tmp/in_N6 > /tmp/in_N4 &

cat /tmp/out_A6 > /tmp/in_C6 &
cat /tmp/out_C6 | tee /tmp/in_A6 > /tmp/in_N6 &
#CONNEXIONS INTERSITES
cat /tmp/out_N6 | tee /tmp/in_C6 | tee /tmp/in_N2 > /tmp/in_N5 &