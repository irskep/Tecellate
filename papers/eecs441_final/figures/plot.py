from matplotlib import pyplot as plt

X = 500

x = [D for D in xrange(1, X)]

t = [1.0 for D in xrange(1, X)]
b = [0.0 for D in xrange(1, X)]

##plt.plot(g)
##plt.plot([1.001 for D in xrange(1, X)], alpha=.0)

#plt.fill_between(x, b,
    #[min(1.0, 1.0/(D/(50.0*1.137))) for D in xrange(1, X)],
    #color='b', alpha=.5)
#plt.fill_between(x, b,
    #[min(1.0, 1.0/(D/(20.0*1.137))) for D in xrange(1, X)],
    #color='r', alpha=.2)
#plt.fill_between(x, b,
    #[min(1.0, 1.0/(D/(10.0*1.137))) for D in xrange(1, X)],
    #color='g', alpha=.2)

#p3 = plt.Rectangle((0, 0), 1, 1, fc="#7f7fff")
#p2 = plt.Rectangle((0, 0), 1, 1, fc="#9865cc")
#p1 = plt.Rectangle((0, 0), 1, 1, fc="#796aa3")
#plt.legend([p1, p2, p3], ["10 units", "20 units", "50 units"], loc=1, title="Broadcast Range")

#plt.title('Probability of a byte not being corrupted.')
#plt.xlabel('Distance')
#plt.ylabel('Probability')
#plt.show()

plt.fill_between(x, b,
    [min(1.0, 1.0/(D/(50.0/(1.137*3)))) for D in xrange(1, X)],
    color='b', alpha=.5)
plt.fill_between(x, b,
    [min(1.0, 1.0/(D/(20.0/(1.137*3)))) for D in xrange(1, X)],
    color='r', alpha=.2)
plt.fill_between(x, b,
    [min(1.0, 1.0/(D/(10.0/(1.137*3)))) for D in xrange(1, X)],
    color='g', alpha=.2)

p1 = plt.Rectangle((0, 0), 1, 1, fc="#7f7fff")
p2 = plt.Rectangle((0, 0), 1, 1, fc="#9865cc")
p3 = plt.Rectangle((0, 0), 1, 1, fc="#796aa3")
plt.legend([p3, p2, p1], ["10 units", "20 units", "50 units"], loc=1, title="Broadcast Range")

plt.title('Probability of a byte being combined.')
plt.xlabel('Distance')
plt.ylabel('Probability')
plt.show()
